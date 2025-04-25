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
      final response = await http.get(Uri.parse('http://127.0.0.1:5566/admin/getallshelterstry'));
      final responseData = json.decode(response.body);

      if (response.statusCode == 200 && responseData['retCode'] == '200') {
        final shelters = responseData['data'] as List;

        final List<dynamic> approved = [];
        final List<dynamic> pending = [];

        for (var shelter in shelters) {
          final shelterStatus = shelter['shelter']?['status'];
          if (shelterStatus == 'active') {
            approved.add(shelter);
          } else {
            pending.add(shelter);
          }
        }

        setState(() {
          approvedShelters = approved;
          pendingShelters = pending;
          isLoading = false;
        });
      } else {
        setState(() {
          isLoading = false;
          errorMessage = 'Failed to load shelters.';
        });
      }
    } catch (e) {
      setState(() {
        isLoading = false;
        errorMessage = 'Error: $e';
      });
    }
  }

  Future<void> _approveShelter(int shelterId) async {
    try {
      final response = await http.post(
        Uri.parse('http://127.0.0.1:5566/admin/approveshelter'),
        body: json.encode({'shelter_id': shelterId}),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Shelter approved successfully!')),
        );
        await _fetchShelters();
      } else {
        throw Exception('Failed to approve shelter');
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error: ${e.toString()}')),
      );
    }
  }

  Widget _buildShelterCard(dynamic shelterData, bool isApproved) {
    final info = shelterData['info'] ?? {};
    final shelter = shelterData['shelter'] ?? {};

    return GestureDetector(
      onTap: () {
        if (!isApproved) {
          _showShelterDetailsDialog(shelterData);
        }
      },
      child: Card(
        elevation: 2,
        margin: const EdgeInsets.symmetric(vertical: 6, horizontal: 12),
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                info['shelter_name'] ?? 'No Name',
                style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 6),
              Text(info['shelter_address'] ?? 'No Address'),
              const SizedBox(height: 6),
              Row(
                children: [
                  const Icon(Icons.phone, size: 16, color: Colors.blue),
                  const SizedBox(width: 4),
                  Text(info['shelter_contact'] ?? 'No Contact'),
                  const Spacer(),
                  Chip(
                    label: Text(
                      shelter['status']?.toUpperCase() ?? 'N/A',
                      style: const TextStyle(color: Colors.white),
                    ),
                    backgroundColor: isApproved ? Colors.green : Colors.orange,
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _showShelterDetailsDialog(dynamic shelterData) {
    final info = shelterData['info'] ?? {};
    final shelter = shelterData['shelter'] ?? {};

    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(15),
          ),
          title: Text(info['shelter_name'] ?? 'No Name'),
          content: Container(
            width: 400, // Adjust width for a bigger rectangular shape
            padding: const EdgeInsets.all(20),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Address: ${info['shelter_address'] ?? 'No Address'}'),
                const SizedBox(height: 10),
                Text('Contact: ${info['shelter_contact'] ?? 'No Contact'}'),
              ],
            ),
          ),
          actions: [
            TextButton(
              onPressed: () {
                Navigator.of(context).pop();
              },
              child: const Text('Cancel'),
            ),
            TextButton(
              onPressed: () async {
                await _approveShelter(shelter['ShelterID']);
                Navigator.of(context).pop();
              },
              child: const Text('Approve'),
            ),
          ],
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    if (isLoading) return const Center(child: CircularProgressIndicator());

    if (errorMessage.isNotEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(errorMessage),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: _fetchShelters,
              child: const Text('Retry'),
            ),
          ],
        ),
      );
    }

    final sheltersToDisplay = showApproved ? approvedShelters : pendingShelters;

    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.all(16),
          child: ToggleButtons(
            isSelected: [!showApproved, showApproved],
            onPressed: (int index) {
              setState(() {
                showApproved = index == 1;
              });
            },
            borderRadius: BorderRadius.circular(20),
            selectedColor: Colors.white,
            fillColor: Colors.blue,
            color: Colors.blue,
            constraints: const BoxConstraints(minHeight: 40, minWidth: 100),
            children: [
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 12),
                child: Text('Pending (${pendingShelters.length})'),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 12),
                child: Text('Approved (${approvedShelters.length})'),
              ),
            ],
          ),
        ),
        Expanded(
          child: ListView.builder(
            itemCount: sheltersToDisplay.length,
            itemBuilder: (context, index) {
              return _buildShelterCard(sheltersToDisplay[index], showApproved);
            },
          ),
        ),
      ],
    );
  }
}
