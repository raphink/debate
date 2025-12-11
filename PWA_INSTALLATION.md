# Mobile PWA Installation Guide

The Theology & Philosophy Debate Generator can be installed as a Progressive Web App (PWA) on your mobile device for a native app-like experience.

## Benefits of Installing

- **Standalone Mode**: Runs without browser chrome (no address bar or navigation buttons)
- **Home Screen Access**: Launch directly from your device's home screen
- **Offline Capability**: Icon and basic assets cached for faster loading
- **Full Screen**: Uses entire screen real estate for better UX

## Installation Instructions

### iOS (Safari)

1. **Open the App**
   - Navigate to the app URL in Safari
   - Example: `https://raphink.github.io/debate`

2. **Access Share Menu**
   - Tap the Share button (box with up arrow) at the bottom of Safari
   - Located in the center of the bottom toolbar

3. **Add to Home Screen**
   - Scroll down in the share sheet
   - Tap "Add to Home Screen"
   - You'll see the app icon and name preview

4. **Confirm Installation**
   - Edit the name if desired (default: "Debate Generator")
   - Tap "Add" in the top right corner
   - The app icon will appear on your home screen

5. **Launch the App**
   - Tap the new icon on your home screen
   - App opens in standalone mode without Safari UI

### Android (Chrome)

1. **Open the App**
   - Navigate to the app URL in Chrome
   - Example: `https://raphink.github.io/debate`

2. **Install Prompt**
   - Chrome may show an automatic "Add to Home Screen" banner
   - If so, tap "Install" or "Add"
   - Skip to step 5

3. **Manual Installation** (if no banner appears)
   - Tap the three-dot menu in the top right
   - Select "Install app" or "Add to Home Screen"

4. **Confirm Installation**
   - Review the app name and icon
   - Tap "Install" or "Add"

5. **Launch the App**
   - Tap the new icon on your home screen or app drawer
   - App opens in standalone mode without Chrome UI

## App Features in PWA Mode

- **Theme Color**: Dark navy blue (#1a1a2e) status bar on Android
- **Orientation**: Portrait-primary (optimized for vertical viewing)
- **Display Mode**: Standalone (no browser UI)
- **Icons**: High-resolution 192x192 and 512x512 icons
- **Scope**: Full access to all app features and pages

## Uninstalling the PWA

### iOS
1. Press and hold the app icon on home screen
2. Tap "Remove App"
3. Select "Delete App"
4. Confirm deletion

### Android
1. Press and hold the app icon
2. Tap "App info" or drag to "Uninstall"
3. Tap "Uninstall"
4. Confirm uninstallation

## Technical Details

- **Manifest**: `/manifest.json` defines app metadata
- **Icons**: SVG source in `/public/icon.svg`
- **Sizes**: 192x192 (standard), 512x512 (high-res), 32x32 (favicon)
- **Purpose**: `any maskable` (adaptive icons on Android)
- **Categories**: Education, Entertainment

## Troubleshooting

**iOS: "Add to Home Screen" option missing**
- Ensure you're using Safari (not Chrome or Firefox on iOS)
- Check iOS version is 14.0 or higher
- Verify the manifest.json loads correctly (check Network tab)

**Android: Install banner doesn't appear**
- Use Chrome browser (version 90+)
- Ensure stable internet connection during first visit
- Check manifest.json is valid (Chrome DevTools > Application > Manifest)
- Try the manual installation method via menu

**App doesn't work offline**
- PWA has no service worker in current version (online-only)
- Ensure internet connection when using the app
- Future versions may add offline capability

## Future Enhancements

- Service worker for true offline support
- Push notifications for debate completion
- Background sync for queued debates
- Share target for receiving debate topics from other apps
