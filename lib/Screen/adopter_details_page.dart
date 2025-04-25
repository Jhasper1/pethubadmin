import 'package:flutter/material.dart';

class AdopterDetailsPage extends StatelessWidget {
  final Map<String, dynamic> adopterData;

  const AdopterDetailsPage({super.key, required this.adopterData});

  @override
  Widget build(BuildContext context) {
    final adopter = adopterData['adopter'];
    final info = adopterData['info'];

    return Scaffold(
      appBar: AppBar(
        title: const Text("Adopter Details"),
        backgroundColor: Colors.teal,
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Card(
          elevation: 4,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
          child: Padding(
            padding: const EdgeInsets.all(20.0),
            child: SingleChildScrollView(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _infoRow("Full Name", "${info['first_name'] ?? ''} ${info['last_name'] ?? ''}"),
                  _infoRow("Age", info['age']?.toString() ?? "N/A"),
                  _infoRow("Sex", info['sex'] ?? "N/A"),
                  _infoRow("Contact Number", info['contact_number'] ?? "N/A"),
                  _infoRow("Email", adopter['email'] ?? "N/A"),
                  _infoRow("Occupation", info['occupation'] ?? "N/A"),
                  _infoRow("Civil Status", info['civil_status'] ?? "N/A"),
                  _infoRow("Social Media", info['social_media'] ?? "N/A"),
                  _infoRow("Address", info['address'] ?? "N/A"),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _infoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            "$label: ",
            style: const TextStyle(
              fontWeight: FontWeight.bold,
              fontSize: 16,
            ),
          ),
          Expanded(
            child: Text(
              value,
              style: const TextStyle(fontSize: 16),
            ),
          ),
        ],
      ),
    );
  }
}
