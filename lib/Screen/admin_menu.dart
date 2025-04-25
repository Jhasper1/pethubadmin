import 'package:flutter/material.dart';

class AdminMenu extends StatefulWidget {
  final bool isCollapsed;

  const AdminMenu({super.key, required this.isCollapsed});

  @override
  State<AdminMenu> createState() => _AdminMenuState();
}

class _AdminMenuState extends State<AdminMenu> {
  bool isExpanded = false;

  void toggleMenu() {
    setState(() {
      isExpanded = !isExpanded;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Flexible(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // When sidebar is collapsed, display only the icon
          if (widget.isCollapsed)
            GestureDetector(
              onTap: toggleMenu,
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                decoration: BoxDecoration(
                  color: Colors.grey.withOpacity(0.2),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const Icon(Icons.person, color: Colors.white70, size: 24),
                  ],
                ),
              ),
            )
          // When sidebar is expanded, display the full name and options
          else
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                GestureDetector(
                  onTap: toggleMenu,
                  child: Container(
                    padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                    decoration: BoxDecoration(
                      color: Colors.grey.withOpacity(0.2),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        const Icon(Icons.person, color: Colors.white70, size: 18),
                        const SizedBox(width: 8),
                        const Text(
                          'Admin Name',
                          style: TextStyle(color: Colors.white70, fontSize: 14),
                        ),
                        Icon(
                          isExpanded
                              ? Icons.arrow_drop_up
                              : Icons.arrow_drop_down,
                          color: Colors.white70,
                        ),
                      ],
                    ),
                  ),
                ),
                // Expanded menu items
                if (isExpanded) ...[
                  _menuItem(Icons.person, "Profile", () {
                    print("Go to profile");
                  }),
                  _menuItem(Icons.settings, "Settings", () {
                    print("Go to settings");
                  }),
                  _menuItem(Icons.logout, "Logout", () {
                    print("Logging out...");
                  }),
                  const SizedBox(height: 10),
                ]
              ],
            ),
        ],
      ),
    );
  }

  Widget _menuItem(IconData icon, String label, VoidCallback onTap) {
    return GestureDetector(
      onTap: () {
        onTap();
        setState(() => isExpanded = false);
      },
      child: Container(
        margin: const EdgeInsets.only(bottom: 6),
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.05),
          borderRadius: BorderRadius.circular(8),
        ),
        child: Row(
          children: [
            Icon(icon, size: 18, color: Colors.white60),
            const SizedBox(width: 8),
            Text(label, style: const TextStyle(color: Colors.white60)),
          ],
        ),
      ),
    );
  }
}
