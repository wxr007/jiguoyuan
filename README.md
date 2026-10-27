# Chiguoyuan / Jiguoyuan

English | **[中文](README_CN.md)**

A standalone, multi-platform short-drama player built with Flutter. Source requests, parsing, downloads and playback all happen on the device; the app does not depend on the legacy projects, a self-hosted remote service or a third-party relay. The native core is Go (loaded over FFI), the player is media_kit / libmpv, and merging and export use FFmpegKit.

Current source version: **0.2.61+70 (unverified development snapshot)**. Per the 2026-09-21 agreement, features come first and centralized testing/regression is on hold; anything not yet validated on a device is marked as pending verification and is not treated as an accepted build.

## Editions and build options

| Build | App name | Available sources |
| --- | --- | --- |
| Default, no arguments | Chiguoyuan (赤果园) | Hongguo only |
| Build script with `--all-sources` | Jiguoyuan (集果园) | All sources |

This is a compile-time option and cannot be switched inside the app. The window title, desktop name, window and distribution file names follow the option, and the interface, source calls and native download scheduler restrict sources accordingly. Both editions share the same application ID and data directory, so they can upgrade over each other; switching editions keeps favorites, watch history, user profiles and download records, and the backup format is compatible.

## Features

| Feature | Details |
| --- | --- |
| Browse and search | Paged category catalog, online search (with suggestions on Hongguo), rankings (Hongguo trending / live-action / animated / AI drama; the full edition adds Huangdou and Huangguo charts), category recommendations |
| Playback | Quality, speed, episode switching, orientation adaptation, auto-play next, automatic route switching, fallback routes, playback error retry; no play or previous/next buttons in the middle of the picture — single tap shows or hides the controls, double tap toggles play/pause, and holding for temporary 3x speed shows only the top hint without the control bar or scrim; an on-device stream server proxies and rewrites HLS playlists while plain MP4 streams pass through |
| Danmaku | Danmaku for Hongguo online playback, toggled per user from a button on the player side |
| Following and history | Follow list, update reminders, watch history, resume from last position, isolated per local user |
| Download and offline | Per-episode downloads at a chosen quality, batch management, resume, Android background service, offline playback, directory migration |
| Merge and export | Episode merging, local finished-media library, Emby export with stable numbering and NFO |
| Users and permissions | Local multi-user support, admin password, per-user source and download permissions, configuration backup and restore |
| LAN interconnect | mDNS discovery, direct HTTPS connection, incremental sync of favorites and progress, push-to-resume without going through a server |
| Video enhancement | Physical output size matching, RAVU, anime CNN, FSRCNNX and more, Windows desktop only (not initialized on mobile or TV) |
| Android TV | Automatic TV detection, forced landscape, remote focus model and option dialogs |
| Theme | Light / dark / follow system |

Source availability, quality and regional restrictions depend on the source sites and your network. The app does not bypass source VIP or any other licensing restrictions.

## Sources

Source names are romanized Chinese platform names; the app interface shows the original Chinese names.

| Source | Browse and playback | Search |
| --- | --- | --- |
| Hongguo | Live-action, animated and AI dramas with episodes | Online search plus official search suggestions |
| Huangdou (full) | List, VIP badges and episodes | Filter already-loaded dramas |
| Juguo (full) | Trending, latest, API categories, detail and playable episodes; signed cookies for playback, preloading and downloads, pending verification | Online source search with load-more |
| Yeguo (full) | API categories, catalog, detail and episodes; route discovery and fallback routes, pending verification | Paged online source search |
| Diguo (full) | Web categories, catalog, detail and episodes; vplayer signature parsing, pending verification | Paged online source search |
| Huangguo Video (full) | List, detail and episodes | Filter already-loaded dramas |
| Huangguo AI (full) | List, categories, detail and episodes; entry connectivity and playback recovery blocked on external conditions, see TODO P1-4 | Filter already-loaded dramas |
| Huangguo Classic (full) | API categories, list, detail and episodes | Filter already-loaded dramas |
| Yaguo (full) | Xingya API categories, catalog, detail and episodes; login token caching, pending verification | Online source search |
| Maoguo (full) | Qimao signed API catalog, detail and episodes, pending verification | Online source search |
| Fanguo (full) | Xifan search API catalog and detail, single-page catalog, pending verification | Online source search |
| Guanguo (full) | Weiguan API catalog and episodes; the detail API returns empty data in testing, pending verification | Online source search |
| Heguo (full) | Hema web `__NEXT_DATA__` catalog, detail and episodes, pending verification | Online source search (`POST /seo/video/6007`) |
| Xingguo (full) | Xingxing serialized API catalog, detail and episodes, pending verification | Online source search (`GET /novel-api/basedata/book/searchBook`) |
| Huaguo (full) | Web categories, catalog, detail and episodes, pending verification | Online source search |
| Niuguo (full) | AES-encrypted API catalog, detail and episodes, pending verification | Online source search |
| Wangguo (full) | Web categories, catalog, detail and episodes, pending verification | Online source search |
| Faguo (full) | Web categories, catalog, detail and episodes, pending verification | Online source search |
| Piguo (full) | Web catalog and detail; `ptt.red` is blocked by a Cloudflare challenge and is not reachable over plain HTTP | Online source search |
| Wuguo (full) | Web categories, catalog, detail and episodes, pending verification | Online source search |
| Qingkong (full) | Paged anime, theatrical animation and tokusatsu categories, detail and episodes; ported source, pending verification | Paged online source search |
| Guipian (full) | Paged horror, TV and anime categories, detail and multi-route episodes; ported source, pending verification | Title matching against the RSS latest entries (the site search API is retired, so matches are limited to the newest batch) |
| Hanxiaoquan (full) | Paged K-drama, Korean film and variety/anime categories, detail and episodes; ported source, pending verification | Paged online source search |

Source management: each entry shows its cached count, paging position, last update time and error, and can be updated individually, advanced by one page, backfilled with detail, stopped, or health-checked; results are persisted and existing caches survive parse failures. Restricted users need a source permission assigned by an admin, and existing users do not automatically gain newly added sources.

## Packages and platform status

| Platform | Status |
| --- | --- |
| Android 8.0+ phones | Source `0.2.61+70`; build scripts and Actions artifacts for ARMv7 / ARM64 / x86_64 are available, installation and runtime verification pending |
| Windows 10/11 x64 | Unzip and run `hongguojian.exe` / `zhenguojian.exe`, keeping all DLLs and the `data` directory; native LAN discovery requires Windows 10 1903+ |
| macOS 12+ | Unzip and run `hongguojian.app` / `zhenguojian.app`; the Go core is a universal arm64 + x86_64 dylib; the bundle is unsigned, so the first launch needs `xattr -cr` to clear the quarantine attribute; build and runtime verification pending |
| Android TV | Shares the Android source with phones; remote control and TV layout pending device verification. The TV banner images still carry the pre-rename wording and must be regenerated on a machine with Pillow via `python scripts/generate_brand_assets.py --font <CJK font path>` |
| iOS 15.1+ | Project, Go core XCFramework, media dependencies and build scripts are in place; the player disables hardware texture acceleration to avoid libmpv exiting; awaiting Xcode build and device verification, no signed IPA yet |

`INSTALL_FAILED_NO_MATCHING_ABIS` means the APK does not match the device architecture; install the package for the right ABI. Android APKs store the native `.so` libraries compressed by default and the system extracts them during installation.

## GitHub Actions

Pushing to `main` / `master` builds only the full-edition (Jiguoyuan) Android APK. Windows, macOS and iOS artifacts are produced only when pushing a `v*` tag or manually running **Build app packages**. The following artifacts are built and kept for 14 days:

| Default edition artifact | Full edition artifact | Contents |
| --- | --- | --- |
| — | `zhenguojian-android` | APKs for three ABIs plus SHA256 |
| `hongguojian-windows` | `zhenguojian-windows` | Complete ZIP plus SHA256 |
| `hongguojian-macos` | `zhenguojian-macos` | Universal macOS `.app` ZIP plus SHA256, unsigned |
| `hongguojian-ios-unsigned` | `zhenguojian-ios-unsigned` | Unsigned `.app` ZIP plus SHA256, not installable as a signed IPA |

The default-edition Android artifact is no longer built; Android ships the full edition only.

The checks job never blocks packaging and runs only on pull requests or manual runs; Flutter and Go test results during the testing pause are informational only. The manual **Publish release** workflow collects artifacts from a chosen run into a GitHub Release; **Convert unsigned iOS app to ipa** packs an unsigned `.app` into an IPA and uploads it to a chosen Release.

Actions minutes on a private repository are billed by runner multiplier (Linux 1x, Windows 2x, macOS 10x), which is why desktop and iOS artifacts are built only on tags or manual runs; only Android builds on every push.

Android release signing is configured through repository secrets:

| Secret | Contents |
| --- | --- |
| `ANDROID_KEYSTORE_BASE64` | Base64 of the JKS file |
| `ANDROID_KEYSTORE_PASSWORD` | Keystore password |
| `ANDROID_KEY_ALIAS` | Key alias |
| `ANDROID_KEY_PASSWORD` | Key password |

The repository already contains the fixed signing key `android/duanju-release.jks` together with `android/key.properties`, so local and CI builds share one signature and their artifacts can upgrade over each other; when the secrets above are configured they take precedence, with the same signature. Replacing the key makes installed versions impossible to upgrade in place, so they must be uninstalled and reinstalled. Only when those files are missing and no secrets are set does the build fall back to debug signing, in which case preview packages from different machines cannot upgrade over each other.

## Development and builds

Flutter `3.47.4`, Dart `3.12+`, Go `1.24.1+`, Python `3.10+`. Android needs JDK 17, SDK 36 and NDK `28.2.13676358`; Windows needs the Visual Studio C++ desktop workload and MinGW-w64 x64; macOS and iOS need macOS, a full Xcode installation and CocoaPods.

Build scripts set `GOPROXY=https://goproxy.cn,direct` and `GOSUMDB=off` for child processes without touching global configuration; environment variables of the same name override them.

~~~sh
python3 scripts/build_android.py            # default edition, add --abi arm64-v8a --cn-mirrors
python3 scripts/build_android.py --all-sources
python3 scripts/build_windows.py            # Windows, default edition
python3 scripts/build_windows.py --all-sources
python3 scripts/build_macos.py              # macOS, default edition (unsigned ZIP)
python3 scripts/build_macos.py --all-sources
python3 scripts/build_ios.py                # unsigned iOS, supports --core-only --simulator
python3 scripts/build_ios.py --export-options /path/to/ExportOptions.plist   # sign and export an IPA yourself
~~~

Windows PowerShell also supports `.\scripts\build_windows.ps1 [-AllSources] [-ChinaMirrors]`.

Builds in mainland China can use `--cn-mirrors` / `-ChinaMirrors`: Flutter and pub use `storage.flutter-io.cn` / `pub.flutter-io.cn` and Android dependencies prefer the Aliyun mirror. This only affects the current build; the original lock file and temporary Gradle configuration are restored afterwards.

For first-time local debugging, build the Go core for the target platform first:

~~~sh
python3 scripts/build_native.py --platform android --abi arm64-v8a
python3 scripts/build_native.py --platform windows
python3 scripts/build_native.py --platform darwin
flutter pub get --enforce-lockfile
flutter run
~~~

When debugging the full edition, pass `--all-sources` to `build_native.py` and run `flutter run --dart-define=ALL_SOURCES=true`; the script sets the Dart constant and the Go build flags together, and the app checks on startup that the two agree.

Artifacts land in `dist/android`, `dist/windows`, `dist/macos` and `dist/ios`. The default edition is prefixed with `hongguojian-` and the full edition with `zhenguojian-`; all come with SHA256SUMS.

## Roadmap

Compared against the feature list of the legacy drama libraries (`../短剧库`, `../果果剧库`). "Both editions" means the default Chiguoyuan and the full-source Jiguoyuan. Android phones and Windows are the primary platforms, Android TV comes next, iOS last.

### Priority fixes: paging, metadata and data boundaries (implemented, awaiting centralized verification)

| ID | Scope | Description | To verify |
| --- | --- | --- | --- |
| P1-R1 | Both | Total catalog and categories persist the real offset, session, device ID, paging signature and last ID in the same snapshot as the catalog | Resume after restart, expired session, overlapping pages, independent category resume |
| P1-R2 | Both | Removed the silent 6000-item and 500-page caps; entries are merged by ID and persisted; the 32 MiB save limit reports a clear error | Item 6001, page 501, restart, cross-category deduplication |
| P1-R4 | Both | Detail fields completed and carried through Go and Dart; detail and playback notifications share the ID list and failed posters | Cross-page write-back, retry of the same failed URL, stale caches |
| P1-R5 | Full, Huangdou | Prefer plain-text JSON before AES decoding, keeping gzip support and long IDs | Plain text, cipher text, malformed payloads |
| P1-R7 | Both | Download index load failures never rewrite a migration; the original file is kept | Corrupt or oversized indexes, disk write failures, interrupted migration |
| P1-R8 | Both | Catalog and source-record write failures are reported back to the page as unsaved with a retry option | File conflicts, space or permission failures, restart |
| P1-R9 | Both | First launch and corrupted configuration are distinguished; access is locked on anomalies with export and admin-verified recovery | Guest or admin corruption, cold start, recovery flow |
| P1-R10 | Both | User settings and records check the save result, roll back to a snapshot on failure and lock when the state cannot be confirmed | Write failures, interrupted recovery, user isolation, restart consistency |

### Feature completion

| ID | Feature | Status and to verify |
| --- | --- | --- |
| P1-4 | Huangguo AI playback recovery | Blocked externally: the catalog and known entries time out on the current development network; real playback needs a reachable entry plus key and media chain |
| P1-6 | Hongguo category recommendations | Implemented, pending verification: recommendation cursor, session and dedupe records persisted, current content kept on failure |
| P1-7 | Sorting and complete metadata | Implemented, pending verification: sorting and filters in the title bar, metadata passed across sources |
| P2-8 | Metadata fields and three-state VIP | Implemented, pending verification: unknown / free / VIP states, unknown never overwrites known, source management supports batched backfill |
| P2-9 | Emby export polish | Implemented, pending verification: stable catalog numbering, Season 00, NFO scope; renames adding episodes, repeated exports, Emby scanning |
| P2-10 | Playback preferences and fullscreen interaction | Implemented, pending verification: quality, speed and auto-play saved per user into backups, fullscreen episode picker and keyboard rules |
| P2-11 | Continue after metadata failures | Implemented, pending verification: retry timing and scan position persisted, 404 cooldown and backoff |
| P3-1 | Network and resource settings | Implemented, pending verification: direct / auto / manual proxy, catalog concurrency and interval, per-platform proxying and switching in flight |
| P3-2 | LAN favorites and progress sync | Implemented, pending verification: auto discovery, direct HTTPS, incremental merge and manual two-way sync; real two-device sync per platform |
| P3-3 | LAN push playback | Implemented, pending verification: reuses the device connection to push a resume request without transferring media URLs, keys or files |
| P3-4 | Video enhancement and super-resolution | Implemented, pending verification: Windows desktop only; quality and sustained performance pending centralized verification |
| P3-5 | Android picture-in-picture | Implemented, pending verification: system PiP integration; TV, Windows and iOS have no system-level PiP |

The legacy project's Docker / HTTP service, browser registration and remote accounts, cross-device server-side records, online STRM signing gateway and scheduled remote Emby sync all depend on an always-on server, so they are out of scope for this app. The corresponding capabilities are on-device users and backups, a native player, local media and Emby file export, plus serverless LAN connections.

## Source sync and recovery

After every completed development, fix or documentation change, run `python3 scripts/finish_task.py --message "what actually changed"` to sync clean source into the sibling `../guoapp` mirror and produce a `果园-YYYYMMDDHHMM.zip` source-only archive (no SDK, dependency directories, caches or build output). Source recovery relies on the version number, the source archive and this README. Sync rules live in `scripts/sync_source.py` and can be verified with `--check`.

## License and components

The player uses [media_kit](https://github.com/media-kit/media-kit) / libmpv, and merging and export use [FFmpegKit min-gpl](https://github.com/sk3llo/ffmpeg_kit_flutter), which bundles FFmpeg, x264 / x265 and other GPL media components; each component keeps its upstream license. FFmpegKit does not transcode during normal playback or downloads; a system FFmpeg is only used for development verification and users do not need to install one.
