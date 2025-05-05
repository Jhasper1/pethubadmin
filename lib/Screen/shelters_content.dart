import 'package:flutter/material.dart';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';

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
  int _rowsPerPage = 10;
  int _currentPage = 0;

  @override
  void initState() {
    super.initState();
    _fetchShelters();
  }

  Future<void> _fetchShelters() async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('auth_token');

    if (token == null || token.isEmpty) {
      throw Exception('No authentication token found');
    }

    setState(() {
      isLoading = true;
      errorMessage = '';
    });

    final url = 'http://127.0.0.1:5566/api/admin/getallshelterstry';
    final headers = {
      "Content-Type": "application/json",
      "Authorization": "Bearer $token",
    };

    try {
      final response = await http.get(Uri.parse(url), headers: headers).timeout(
            const Duration(seconds: 10),
          );

      if (response.statusCode == 200) {
        final responseData = json.decode(response.body);
        debugPrint("API Response: $responseData");

        if (responseData['retCode'] == '200') {
          final shelters = responseData['data'] as List;
          _categorizeShelters(shelters);
        } else if (responseData['retCode'] == '401') {
          _handleUnauthorized();
        } else {
          throw Exception(
              'Failed to load shelters. Message: ${responseData['message'] ?? 'Unknown error'}');
        }
      } else {
        throw Exception('Failed to fetch data: ${response.statusCode}');
      }
    } catch (e) {
      debugPrint('Error fetching shelters: $e');
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

  void _handleUnauthorized() {
    setState(() {
      errorMessage = 'Session expired. Please log in again.';
    });
    Future.delayed(const Duration(seconds: 2), () {
      Navigator.pushReplacementNamed(context, '/login');
    });
  }

  void _categorizeShelters(List<dynamic> shelters) {
    final List<dynamic> approved = [];
    final List<dynamic> pending = [];

    for (var shelter in shelters) {
      if (_isShelterApproved(shelter)) {
        approved.add(shelter);
      } else {
        pending.add(shelter);
      }
    }

    setState(() {
      approvedShelters = approved;
      pendingShelters = pending;
    });
  }

  bool _isShelterApproved(dynamic shelter) {
    final status = (shelter['reg_status'] ??
            shelter['shelter']?['reg_status'] ??
            shelter['shelter_id']?['reg_status'] ??
            'pending')
        .toString()
        .toLowerCase();
    return status == 'approved';
  }

  int? _extractShelterId(dynamic shelterData) {
    if (shelterData['shelter_id'] is int) return shelterData['shelter_id'];
    if (shelterData['shelter'] is Map &&
        shelterData['shelter']['shelter_id'] is int) {
      return shelterData['shelter']['shelter_id'];
    }
    if (shelterData['id'] is int) return shelterData['id'];
    return null;
  }

  Future<void> _approveShelter(dynamic shelterData) async {
    try {
      final shelterId = _extractShelterId(shelterData);
      if (shelterId == null) throw Exception('Invalid shelter ID');

      final prefs = await SharedPreferences.getInstance();
      final token = prefs.getString('auth_token');
      if (token == null) throw Exception('No authentication token found');

      debugPrint('Approving shelter ID: $shelterId');

      final response = await http.put(
        Uri.parse(
            'http://127.0.0.1:5566/api/admin/shelters/$shelterId/approve'),
        headers: {
          'Content-Type': 'application/json',
          "Authorization": "Bearer $token",
        },
      ).timeout(const Duration(seconds: 10));

      if (response.statusCode == 200) {
        _handleApprovalSuccess(shelterData, shelterId);
        _showSuccessMessage('Shelter approved successfully!');
        Navigator.of(context).pop();
      } else {
        final error = json.decode(response.body);
        throw Exception(error['message'] ??
            'Approval failed with status ${response.statusCode}');
      }
    } catch (e) {
      debugPrint('Error in _approveShelter: $e');
      _showErrorMessage(
          'Failed to approve: ${e.toString().replaceAll('Exception: ', '')}');
    }
  }

  void _handleApprovalSuccess(dynamic shelterData, int shelterId) {
    final index =
        pendingShelters.indexWhere((s) => _extractShelterId(s) == shelterId);
    if (index == -1) return;

    final approvedShelter = _updateShelterStatus(shelterData);

    setState(() {
      pendingShelters.removeAt(index);
      approvedShelters.add(approvedShelter);
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
      SnackBar(
        content: Text(message),
        backgroundColor: Colors.green,
        duration: const Duration(seconds: 3),
      ),
    );
  }

  void _showErrorMessage(String message) {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: Colors.red,
        duration: const Duration(seconds: 3),
      ),
    );
  }

  void _showShelterDialog(dynamic shelterData) {
    final info = shelterData['info'] ?? {};
    final status = _isShelterApproved(shelterData) ? 'Approved' : 'Pending';

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Shelter Details ($status)'),
        content: SizedBox(
          width: MediaQuery.of(context).size.width * 0.7,
          child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _buildDetailRow(
                    'ID', _extractShelterId(shelterData)?.toString() ?? 'N/A'),
                _buildDetailRow('Name', info['shelter_name']),
                _buildDetailRow('Address', info['shelter_address']),
                _buildDetailRow('Contact', info['shelter_contact']),
                _buildDetailRow('Email', info['shelter_email']),
                _buildDetailRow('Status', status),
                if (info['description'] != null)
                  _buildDetailRow('Description', info['description']),
              ],
            ),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Close'),
          ),
          if (!_isShelterApproved(shelterData))
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
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 100,
            child: Text(
              '$label:',
              style: const TextStyle(fontWeight: FontWeight.bold),
            ),
          ),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              value ?? 'Not provided',
              softWrap: true,
            ),
          ),
        ],
      ),
    );
  }

  List<dynamic> _getPaginatedShelters() {
    final shelters = showApproved ? approvedShelters : pendingShelters;
    final startIndex = _currentPage * _rowsPerPage;
    final endIndex = startIndex + _rowsPerPage;
    return shelters.sublist(
        startIndex, endIndex > shelters.length ? shelters.length : endIndex);
  }

  @override
  Widget build(BuildContext context) {
    if (isLoading) {
      return const Center(
        child: CircularProgressIndicator(),
      );
    }

    if (errorMessage.isNotEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              errorMessage,
              style: const TextStyle(color: Colors.red, fontSize: 16),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 20),
            ElevatedButton(
              onPressed: _fetchShelters,
              style: ElevatedButton.styleFrom(
                padding:
                    const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
              ),
              child: const Text('Retry'),
            ),
          ],
        ),
      );
    }

    final shelters = showApproved ? approvedShelters : pendingShelters;
    final paginatedShelters = _getPaginatedShelters();

    return Scaffold(
      appBar: AppBar(
        title: const Text('Shelter Management'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: _fetchShelters,
            tooltip: 'Refresh',
          ),
        ],
      ),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  showApproved
                      ? 'Approved Shelters (${approvedShelters.length})'
                      : 'Pending Shelters (${pendingShelters.length})',
                  style: const TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                Row(
                  children: [
                    const Text('Show Approved'),
                    Switch(
                      value: showApproved,
                      onChanged: (value) {
                        setState(() {
                          showApproved = value;
                          _currentPage = 0;
                        });
                      },
                      activeColor: Colors.blue,
                    ),
                  ],
                ),
              ],
            ),
          ),
          Expanded(
            child: Column(
              children: [
                Expanded(
                  child: SingleChildScrollView(
                    scrollDirection: Axis.horizontal,
                    child: SingleChildScrollView(
                      scrollDirection: Axis.vertical,
                      child: DataTable(
                        columns: const [
                          DataColumn(
                              label: Text('ID',
                                  style:
                                      TextStyle(fontWeight: FontWeight.bold))),
                          DataColumn(
                              label: Text('Name',
                                  style:
                                      TextStyle(fontWeight: FontWeight.bold))),
                          DataColumn(
                              label: Text('Contact',
                                  style:
                                      TextStyle(fontWeight: FontWeight.bold))),
                          DataColumn(
                              label: Text('Email',
                                  style:
                                      TextStyle(fontWeight: FontWeight.bold))),
                          DataColumn(
                              label: Text('Status',
                                  style:
                                      TextStyle(fontWeight: FontWeight.bold))),
                          DataColumn(
                              label: Text('Actions',
                                  style:
                                      TextStyle(fontWeight: FontWeight.bold))),
                        ],
                        rows: paginatedShelters.map((shelter) {
                          final info = shelter['info'] ?? {};
                          final shelterId =
                              _extractShelterId(shelter)?.toString() ?? 'N/A';
                          final status = _isShelterApproved(shelter)
                              ? 'Approved'
                              : 'Pending';

                          return DataRow(
                            cells: [
                              DataCell(Text(shelterId)),
                              DataCell(
                                ConstrainedBox(
                                  constraints:
                                      const BoxConstraints(maxWidth: 150),
                                  child: Text(
                                    info['shelter_name'] ?? 'Unnamed Shelter',
                                    overflow: TextOverflow.ellipsis,
                                  ),
                                ),
                              ),
                              DataCell(Text(
                                  info['shelter_contact'] ?? 'Not provided')),
                              DataCell(Text(
                                  info['shelter_email'] ?? 'Not provided')),
                              DataCell(
                                Chip(
                                  label: Text(
                                    status,
                                    style: TextStyle(
                                      color: status == 'Approved'
                                          ? Colors.green
                                          : Colors.orange,
                                    ),
                                  ),
                                  backgroundColor: status == 'Approved'
                                      ? Colors.green[50]
                                      : Colors.orange[50],
                                ),
                              ),
                              DataCell(
                                Row(
                                  children: [
                                    IconButton(
                                      icon: const Icon(Icons.visibility,
                                          size: 20),
                                      onPressed: () =>
                                          _showShelterDialog(shelter),
                                      tooltip: 'View Details',
                                    ),
                                  ],
                                ),
                              ),
                            ],
                          );
                        }).toList(),
                        headingRowColor:
                            MaterialStateProperty.resolveWith<Color>(
                          (Set<MaterialState> states) => Colors.grey[200]!,
                        ),
                        dataRowHeight: 60,
                        headingRowHeight: 50,
                        horizontalMargin: 12,
                        columnSpacing: 20,
                        showCheckboxColumn: false,
                      ),
                    ),
                  ),
                ),
                if (shelters.length > _rowsPerPage)
                  Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.end,
                      children: [
                        DropdownButton<int>(
                          value: _rowsPerPage,
                          items: [5, 10, 25, 50].map((int value) {
                            return DropdownMenuItem<int>(
                              value: value,
                              child: Text('$value per page'),
                            );
                          }).toList(),
                          onChanged: (value) {
                            if (value != null) {
                              setState(() {
                                _rowsPerPage = value;
                                _currentPage = 0;
                              });
                            }
                          },
                        ),
                        const SizedBox(width: 16),
                        Text(
                          '${_currentPage * _rowsPerPage + 1}-${(_currentPage + 1) * _rowsPerPage > shelters.length ? shelters.length : (_currentPage + 1) * _rowsPerPage} of ${shelters.length}',
                        ),
                        const SizedBox(width: 16),
                        IconButton(
                          icon: const Icon(Icons.chevron_left),
                          onPressed: _currentPage > 0
                              ? () {
                                  setState(() {
                                    _currentPage--;
                                  });
                                }
                              : null,
                        ),
                        IconButton(
                          icon: const Icon(Icons.chevron_right),
                          onPressed: (_currentPage + 1) * _rowsPerPage <
                                  shelters.length
                              ? () {
                                  setState(() {
                                    _currentPage++;
                                  });
                                }
                              : null,
                        ),
                      ],
                    ),
                  ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
