import 'package:flutter/material.dart';
import 'dart:convert';
import 'package:http/http.dart' as http;

class SheltersContent extends StatefulWidget {
  const SheltersContent({Key? key}) : super(key: key);

  @override
  _SheltersContentState createState() => _SheltersContentState();
}

class _SheltersContentState extends State<SheltersContent> {
  List<dynamic> pendingShelters = [];
  List<dynamic> approvedShelters = [];
  bool isLoading = true;
  String errorMessage = '';
  bool showApproved = false;
  int currentIndex = 0;

  @override
  void initState() {
    super.initState();
    _fetchShelters();
  }

  Future<void> _fetchShelters() async {
    setState(() {
      isLoading = true;
      errorMessage = '';
    });

    try {
      final response = await http.get(
        Uri.parse('http://127.0.0.1:5566/admin/getallshelterstry'),
        headers: {'Content-Type': 'application/json'},
      );

      final responseData = json.decode(response.body);

      if (response.statusCode == 200 && responseData['retCode'] == '200') {
        final shelters = responseData['data'] as List;
        _categorizeShelters(shelters);
      } else {
        throw Exception(
            'Failed to load shelters. Status: ${response.statusCode}, Message: ${responseData['message']}');
      }
    } catch (e) {
      setState(() {
        errorMessage = 'Error: ${e.toString()}';
      });
    } finally {
      if (mounted) {
        setState(() {
          isLoading = false;
        });
      }
    }
  }

  void _categorizeShelters(List<dynamic> shelters) {
    final List<dynamic> approved = [];
    final List<dynamic> pending = [];

    for (var shelter in shelters) {
      final status = _getShelterStatus(shelter);
      if (status == 'approved') {
        approved.add(shelter);
      } else {
        pending.add(shelter);
      }
    }

    setState(() {
      approvedShelters = approved;
      pendingShelters = pending;
      currentIndex = pending.isNotEmpty ? 0 : -1;
    });
  }

  String _getShelterStatus(dynamic shelter) {
    return (shelter['reg_status'] ??
        shelter['shelter']?['reg_status'] ??
        shelter['shelter_id']?['reg_status'] ??
        'pending')
        .toString()
        .toLowerCase();
  }

  // Update to correctly fetch shelter_id
  int? _extractShelterId(dynamic shelterData) {
    if (shelterData['shelter_id'] is int) return shelterData['shelter_id'];
    if (shelterData['shelter'] is Map && shelterData['shelter']['shelter_id'] is int) {
      return shelterData['shelter']['shelter_id'];
    }
    return null;
  }

  Future<void> _approveShelter(dynamic shelterData) async {
    try {
      final shelterId = _extractShelterId(shelterData);
      if (shelterId == null) throw Exception('Invalid shelter ID');

      debugPrint('Approving shelter ID: $shelterId');

      final response = await http.put(
        Uri.parse('http://127.0.0.1:5566/admin/shelters/$shelterId/approve'),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        _handleApprovalSuccess(shelterData, shelterId);
        _showSuccessMessage('Shelter approved successfully!');
        Navigator.of(context).pop(); // Close the dialog
      } else {
        final error = json.decode(response.body);
        throw Exception(error['message'] ?? 'Approval failed with status ${response.statusCode}');
      }
    } catch (e) {
      debugPrint('Approval error: $e');
      _showErrorMessage('Error: ${e.toString()}');
    }
  }

  void _handleApprovalSuccess(dynamic shelterData, int shelterId) {
    final index = pendingShelters.indexWhere((s) => _extractShelterId(s) == shelterId);
    if (index == -1) return;

    final approvedShelter = _updateShelterStatus(shelterData);

    setState(() {
      pendingShelters.removeAt(index);
      approvedShelters.add(approvedShelter);
      currentIndex = pendingShelters.isEmpty ? -1 : (currentIndex % pendingShelters.length);
    });
  }

  dynamic _updateShelterStatus(dynamic shelter) {
    final updated = Map<String, dynamic>.from(shelter);

    void deepUpdate(Map<String, dynamic> data) {
      data['reg_status'] = 'approved';
      if (data['shelter_id'] is Map) deepUpdate(data['shelter_id']);
      if (data['shelter'] is Map) deepUpdate(data['shelter']);
    }

    deepUpdate(updated);
    return updated;
  }

  void _showSuccessMessage(String message) {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message), backgroundColor: Colors.green),
    );
  }

  void _showErrorMessage(String message) {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message), backgroundColor: Colors.red),
    );
  }

  void _showShelterDialog(dynamic shelterData) {
    final info = shelterData['info'] ?? {};
    final status = _getShelterStatus(shelterData);

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Shelter Details (${status.toUpperCase()})'),
        content: Container(
          width: 700, // Adjust the width here (e.g., 300)
          child: SingleChildScrollView(
            child: ListBody(
              children: [
                _buildDetailRow('Name', info['shelter_name']),
                _buildDetailRow('Address', info['shelter_address']),
                _buildDetailRow('Contact', info['shelter_contact']),
                _buildDetailRow('Email', info['shelter_email']),
                _buildDetailRow('Status', status),
              ],
            ),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Close'),
          ),
          if (status != 'approved')
            ElevatedButton(
              onPressed: () => _approveShelter(shelterData),
              child: const Text('Approve'),
            ),
        ],
      ),
    );


  }

  Widget _buildDetailRow(String label, String? value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text('$label: ', style: const TextStyle(fontWeight: FontWeight.bold)),
          Expanded(child: Text(value ?? 'N/A')),
        ],
      ),
    );
  }

  Widget _buildShelterCard(dynamic shelterData, {bool isApproved = false}) {
    final info = shelterData['info'] ?? {};
    final status = _getShelterStatus(shelterData);

    return Card(
      elevation: 2,
      margin: const EdgeInsets.symmetric(vertical: 8, horizontal: 16),
      child: InkWell(
        onTap: () => _showShelterDialog(shelterData),
        borderRadius: BorderRadius.circular(8),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    info['shelter_name'] ?? 'Unnamed Shelter',
                    style: const TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: status == 'approved' ? Colors.green[100] : Colors.orange[100],
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      status.toUpperCase(),
                      style: TextStyle(
                        color: status == 'approved' ? Colors.green[800] : Colors.orange[800],
                        fontSize: 12,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 8),
              Text(
                info['shelter_address'] ?? 'No address provided',
                style: TextStyle(color: Colors.grey[600]),
              ),
              if (info['shelter_contact'] != null) ...[
                const SizedBox(height: 4),
                Text(
                  'Contact: ${info['shelter_contact']}',
                  style: TextStyle(color: Colors.grey[600]),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }

  void _showPrevious() {
    if (pendingShelters.isEmpty) return;
    setState(() {
      currentIndex = (currentIndex - 1 + pendingShelters.length) % pendingShelters.length;
    });
  }

  void _showNext() {
    if (pendingShelters.isEmpty) return;
    setState(() {
      currentIndex = (currentIndex + 1) % pendingShelters.length;
    });
  }

  @override
  Widget build(BuildContext context) {
    if (isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (errorMessage.isNotEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(errorMessage, style: const TextStyle(color: Colors.red)),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: _fetchShelters,
              child: const Text('Retry'),
            ),
          ],
        ),
      );
    }

    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.all(16),
          child: ToggleButtons(
            isSelected: [!showApproved, showApproved],
            onPressed: (int index) {
              setState(() {
                showApproved = index == 1;
                currentIndex = 0;
              });
            },
            borderRadius: BorderRadius.circular(8),
            selectedColor: Colors.white,
            fillColor: Colors.blue,
            color: Colors.blue,
            constraints: const BoxConstraints(minHeight: 40, minWidth: 120),
            children: [
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16),
                child: Text('Pending (${pendingShelters.length})'),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16),
                child: Text('Approved (${approvedShelters.length})'),
              ),
            ],
          ),
        ),
        Expanded(
          child: showApproved
              ? _buildApprovedList()
              : _buildPendingCarousel(),
        ),
      ],
    );
  }

  Widget _buildApprovedList() {
    if (approvedShelters.isEmpty) {
      return const Center(child: Text('No approved shelters found.'));
    }

    return RefreshIndicator(
      onRefresh: _fetchShelters,
      child: ListView.builder(
        padding: const EdgeInsets.only(bottom: 16),
        itemCount: approvedShelters.length,
        itemBuilder: (context, index) {
          return _buildShelterCard(approvedShelters[index], isApproved: true);
        },
      ),
    );
  }

  Widget _buildPendingCarousel() {
    if (pendingShelters.isEmpty) {
      return const Center(child: Text('No pending shelters found.'));
    }

    return RefreshIndicator(
      onRefresh: _fetchShelters,
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          _buildShelterCard(pendingShelters[currentIndex]),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              IconButton(
                onPressed: _showPrevious,
                icon: const Icon(Icons.arrow_left),
              ),
              IconButton(
                onPressed: _showNext,
                icon: const Icon(Icons.arrow_right),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
