import 'package:flutter/material.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

class DashboardContent extends StatefulWidget {
  const DashboardContent({super.key});

  @override
  State<DashboardContent> createState() => _DashboardContentState();
}

class _DashboardContentState extends State<DashboardContent> {
  int shelterCount = 0;
  int adopterCount = 0;
  int petCount = 0;
  int adoptedPetCount = 0;
  int approvedShelters = 0;
  int pendingShelters = 0;
  bool isLoading = true;
  String errorMessage = '';

  @override
  void initState() {
    super.initState();
    _fetchDashboardData();
  }

  Future<void> _fetchDashboardData() async {
    const baseUrl = 'http://127.0.0.1:5566';

    try {
      final responses = await Future.wait([
        http.get(Uri.parse('$baseUrl/admin/shelters/count')),
        http.get(Uri.parse('$baseUrl/admin/adopters/count')),
        http.get(Uri.parse('$baseUrl/admin/pets/count')),
        http.get(Uri.parse('$baseUrl/admin/adoptedpets/count')),
        http.get(Uri.parse('$baseUrl/admin/pendingshelters/count')),
        http.get(Uri.parse('$baseUrl/admin/approvedshelters/count')),
      ]);

      setState(() {
        if (responses[0].statusCode == 200) {
          shelterCount = json.decode(responses[0].body)['count'] ?? 0;
        }
        if (responses[1].statusCode == 200) {
          adopterCount = json.decode(responses[1].body)['count'] ?? 0;
        }
        if (responses[2].statusCode == 200) {
          petCount = json.decode(responses[2].body)['count'] ?? 0;
        }
        if (responses[3].statusCode == 200) {
          adoptedPetCount = json.decode(responses[3].body)['count'] ?? 0;
        }
        if (responses[4].statusCode == 200) {
          pendingShelters = json.decode(responses[4].body)['count'] ?? 0;
        }
        if (responses[5].statusCode == 200) {
          approvedShelters = json.decode(responses[5].body)['count'] ?? 0;
        }

        isLoading = false;
      });
    } catch (e) {
      setState(() {
        errorMessage = 'Failed to load data: ${e.toString()}';
        isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    if (isLoading) return const Center(child: CircularProgressIndicator());
    if (errorMessage.isNotEmpty) return Center(child: Text(errorMessage));

    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            'Pethub Dashboard',
            style: TextStyle(
              fontSize: 28,
              fontWeight: FontWeight.bold,
              color: Colors.blue,
            ),
          ),
          const SizedBox(height: 20),

          Row(
            children: [
              Expanded(child: _buildStatCard('Shelters', shelterCount, Icons.home)),
              const SizedBox(width: 16),
              Expanded(child: _buildStatCard('Adopters', adopterCount, Icons.people)),
              const SizedBox(width: 16),
              Expanded(child: _buildStatCard('Pets', petCount, Icons.pets)),
              const SizedBox(width: 16),
              Expanded(child: _buildStatCard('Adopted Pets', adoptedPetCount, Icons.favorite)),
            ],
          ),
          const SizedBox(height: 30),

          const Text(
            'Shelter Status',
            style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 10),

          if (shelterCount > 0)
            SizedBox(
              height: 300,
              child: PieChart(
                PieChartData(
                  sections: [
                    PieChartSectionData(
                      value: approvedShelters.toDouble(),
                      color: Colors.green,
                      title:
                      '${((approvedShelters / shelterCount) * 100).toStringAsFixed(1)}%',
                      radius: 60,
                      titleStyle: const TextStyle(
                          color: Colors.white, fontWeight: FontWeight.bold),
                    ),
                    PieChartSectionData(
                      value: pendingShelters.toDouble(),
                      color: Colors.orange,
                      title:
                      '${((pendingShelters / shelterCount) * 100).toStringAsFixed(1)}%',
                      radius: 60,
                      titleStyle: const TextStyle(
                          color: Colors.white, fontWeight: FontWeight.bold),
                    ),
                  ],
                  sectionsSpace: 2,
                  centerSpaceRadius: 60,
                ),
              ),
            ),

          if (shelterCount > 0)
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                _buildLegendItem('Approved', Colors.green),
                const SizedBox(width: 20),
                _buildLegendItem('Pending', Colors.orange),
              ],
            ),
        ],
      ),
    );
  }

  Widget _buildStatCard(String title, int count, IconData icon) {
    return Card(
      elevation: 4,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start, // Align everything to the left
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              count.toString(),
              style: const TextStyle(
                fontSize: 32,
                fontWeight: FontWeight.bold,
                color: Colors.blue,
              ),
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Icon(icon, size: 24, color: Colors.blue),
                const SizedBox(width: 8),
                Text(
                  title,
                  style: const TextStyle(fontSize: 16),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }



  Widget _buildLegendItem(String text, Color color) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Container(width: 16, height: 16, color: color),
        const SizedBox(width: 8),
        Text(text),
      ],
    );
  }
}
