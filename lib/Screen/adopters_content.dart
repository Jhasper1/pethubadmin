import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'adopter_details_page.dart';

class AdoptersPage extends StatefulWidget {
  const AdoptersPage({super.key});

  @override
  State<AdoptersPage> createState() => _AdoptersPageState();
}

class _AdoptersPageState extends State<AdoptersPage> {
  bool isLoading = true;
  List<Map<String, dynamic>> adopters = [];

  @override
  void initState() {
    super.initState();
    _fetchAdopters();
  }

  Future<void> _fetchAdopters() async {
    final response =
    await http.get(Uri.parse('http://127.0.0.1:5566/admin/getalladopters'));

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      setState(() {
        isLoading = false;
        adopters = List<Map<String, dynamic>>.from(data['data']);
      });
    } else {
      setState(() {
        isLoading = false;
      });
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Failed to load adopters')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.white70,
        elevation: 1,
        title: const Align(
          alignment: Alignment.centerLeft,
          child: Text(
            'Adopters',
            style: TextStyle(color: Colors.black),
          ),
        ),
        iconTheme: const IconThemeData(color: Colors.black),
      ),
      body: isLoading
          ? const Center(child: CircularProgressIndicator())
          : Padding(
        padding: const EdgeInsets.all(16.0),
        child: ListView.builder(
          itemCount: adopters.length,
          itemBuilder: (context, index) {
            final adopter = adopters[index];
            final adopterInfo = adopter['info'];
            final adopterAccount = adopter['adopter'];

            return Card(
              elevation: 5,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
              margin: const EdgeInsets.symmetric(vertical: 8),
              child: ListTile(
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) =>
                          AdopterDetailsPage(adopterData: adopter),
                    ),
                  );
                },
                contentPadding: const EdgeInsets.all(15),
                leading: CircleAvatar(
                  radius: 30,
                  backgroundColor: const Color(0xFF1E1E2C),
                  child: const Icon(
                    Icons.person,
                    size: 40,
                    color: Colors.white,
                  ),
                ),
                title: Text(
                  '${adopterInfo['first_name'] ?? ''} ${adopterInfo['last_name'] ?? ''}',
                  style: const TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                subtitle: Padding(
                  padding: const EdgeInsets.only(top: 5),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Email: ${adopterAccount['email'] ?? 'N/A'}',
                        style: TextStyle(color: Colors.grey[600]),
                      ),
                      Text(
                        'Phone: ${adopterInfo['phone'] ?? 'N/A'}',
                        style: TextStyle(color: Colors.grey[600]),
                      ),
                      Text(
                        'Address: ${adopterInfo['address'] ?? 'N/A'}',
                        style: TextStyle(color: Colors.grey[600]),
                      ),
                    ],
                  ),
                ),
                // trailing removed
              ),
            );
          },
        ),
      ),
    );
  }
}
