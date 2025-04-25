import 'package:flutter/material.dart';
import 'dashboard_content.dart';
import 'shelters_content.dart';
import 'adopters_content.dart';
import 'settings_content.dart';

class DashboardScreen extends StatefulWidget {
  const DashboardScreen({super.key});

  @override
  State<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends State<DashboardScreen> {
  final GlobalKey<ScaffoldState> _scaffoldKey = GlobalKey<ScaffoldState>();
  String _selectedItem = 'Dashboard';
  final String adminName = 'Admin User';
  bool _isHovered = false;

  Widget _getSelectedContent() {
    switch (_selectedItem) {
      case 'Shelters':
        return const SheltersContent();
      case 'Adopters':
        return const AdoptersPage();
      case 'Settings':
        return const SettingsContent();
      default:
        return const DashboardContent();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      key: _scaffoldKey,
      body: Row(
        children: [
          // Sidebar
          MouseRegion(
            onEnter: (_) => setState(() => _isHovered = true),
            onExit: (_) => setState(() => _isHovered = false),
            child: AnimatedContainer(
              duration: const Duration(milliseconds: 200),
              width: _isHovered ? 200 : 60,
              height: double.infinity,
              color: Color(0xff0d0d27),
              child: Stack(
                children: [
                  Column(
                    children: [
                      // Logo
                      Container(
                        height: 80,
                        padding: const EdgeInsets.all(10),
                        child: Row(
                          children: [
                            const Icon(Icons.pets, color: Colors.white, size: 30),
                            if (_isHovered)
                              const Padding(
                                padding: EdgeInsets.only(left: 10),
                                child: Text(
                                  'PetAdopt',
                                  style: TextStyle(
                                    color: Colors.white,
                                    fontSize: 20,
                                    fontWeight: FontWeight.bold,
                                  ),
                                ),
                              ),
                          ],
                        ),
                      ),
                      const Divider(color: Colors.white54, height: 1),

                      // Menu Items
                      Expanded(
                        child: ListView(
                          padding: EdgeInsets.zero,
                          children: [
                            _buildMenuItem('Dashboard', Icons.dashboard),
                            _buildMenuItem('Shelters', Icons.home_work),
                            _buildMenuItem('Adopters', Icons.people),
                            _buildMenuItem('Settings', Icons.settings),
                          ],
                        ),
                      ),
                    ],
                  ),

                  // Admin Info
                  Positioned(
                    bottom: 0,
                    left: 0,
                    right: 0,
                    child: Container(
                      padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 10),
                      decoration: const BoxDecoration(
                        border: Border(top: BorderSide(color: Colors.white54)),
                      ),
                      child: Row(
                        children: [
                          const CircleAvatar(
                            radius: 15,
                            backgroundColor: Colors.white,
                            child: Icon(Icons.person, size: 18, color: Colors.blue),
                          ),
                          if (_isHovered)
                            Padding(
                              padding: const EdgeInsets.only(left: 10),
                              child: Text(
                                adminName,
                                style: const TextStyle(color: Colors.white),
                              ),
                            ),
                        ],
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),

          // Main Content Area
          Expanded(
            child: Container(
              color: Colors.grey[100],
              child: _getSelectedContent(),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildMenuItem(String title, IconData icon) {
    bool isSelected = _selectedItem == title;
    return MouseRegion(
      cursor: SystemMouseCursors.click,
      child: GestureDetector(
        onTap: () {
          setState(() {
            _selectedItem = title;
          });
        },
        child: Container(
          color: isSelected ? Color(0xff3b3b40) : Colors.transparent,
          padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 10),
          child: Row(
            children: [
              Icon(icon, color: Colors.white),
              if (_isHovered)
                Padding(
                  padding: const EdgeInsets.only(left: 10),
                  child: Text(
                    title,
                    style: const TextStyle(color: Colors.white),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}
