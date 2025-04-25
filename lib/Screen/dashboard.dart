import 'package:flutter/material.dart';
import 'dashboard_content.dart';
import 'shelters_content.dart';
import 'adopters_content.dart';
import 'settings_content.dart';
import 'admin_menu.dart';

class DashboardScreen extends StatefulWidget {
  const DashboardScreen({super.key});

  @override
  State<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends State<DashboardScreen> {
  int selectedIndex = 0;
  bool isSidebarExpanded = true;

  final List<String> pageTitles = [
    "Dashboard",
    "Shelters",
    "Adopters",
    "Settings",
  ];

  final List<IconData> pageIcons = [
    Icons.dashboard,
    Icons.pets,
    Icons.people,
    Icons.settings,
  ];

  final List<Widget> pages = [
    DashboardContent(),
    SheltersContent(),
    AdoptersPage(),
    SettingsContent(),
  ];

  void _handleMouseEnter(_) {
    setState(() => isSidebarExpanded = false);
  }

  void _handleMouseExit(_) {
    setState(() => isSidebarExpanded = true);
  }

  static const sidebarColor = Color(0xFF1E1E2C);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Row(
        children: [
          // Sidebar
          AnimatedContainer(
            duration: const Duration(milliseconds: 200),
            width: isSidebarExpanded ? 250 : 70,
            color: sidebarColor,
            child: Column(
              children: [
                const SizedBox(height: 24),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 16),
                  child: Image.asset(
                    'assets/images/transparent.png',
                    width: isSidebarExpanded ? 80 : 40,
                  ),
                ),
                const SizedBox(height: 40),
                for (int i = 0; i < pageTitles.length; i++) ...[
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 8),
                    child: ListTile(
                      leading: Icon(pageIcons[i], color: Colors.white70),
                      title: isSidebarExpanded
                          ? Text(
                        pageTitles[i],
                        style: const TextStyle(color: Colors.white70),
                      )
                          : null,
                      tileColor: selectedIndex == i
                          ? Colors.white.withOpacity(0.1)
                          : null,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(8),
                      ),
                      onTap: () => setState(() => selectedIndex = i),
                    ),
                  ),
                ],
                const Spacer(),
                Padding(
                  padding: const EdgeInsets.all(16),
                  child: AdminMenu(isCollapsed: !isSidebarExpanded),
                ),
              ],
            ),
          ),

          // Hover-detecting content area
          Expanded(
            child: MouseRegion(
              onEnter: _handleMouseEnter,
              onExit: _handleMouseExit,
              child: Container(
                color: Colors.grey[50],
                child: Padding(
                  padding: const EdgeInsets.all(24),
                  child: Container(
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(12),
                      boxShadow: [
                        BoxShadow(
                          color: Colors.black.withOpacity(0.05),
                          blurRadius: 10,
                          offset: const Offset(0, 4),
                        ),
                      ],
                    ),
                    child: pages[selectedIndex],
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

