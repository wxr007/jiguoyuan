import argparse
import hashlib
import os
import platform
import shutil
import subprocess
import sys
import time
import urllib.request
from pathlib import Path

from build_mirrors import china_mirror_environment, mirrored_pub_lockfile
from app_build import BuildVariant, add_variant_argument

LIBMPV_VERSION = 'v1.1.7'
LIBMPV_URL = 'https://github.com/media-kit/libmpv-android-video-build/releases/download/{version}/default-{abi}.jar'
LIBMPV_JARS = (
    ('arm64-v8a', '83df25b61193af8fa815e373143ac9af'),
    ('armeabi-v7a', '22e21526fefc0a2b8f17adbec9f57590'),
    ('x86_64', '6fa26bf0459b11f1c0b0dbc29e5b940d'),
    ('x86', '0d742b756dc9d1fcd84ea271d8b68f32'),
)


def _plugin_version_key(path):
    parts = path.name.split('-')[-1].split('.')
    return tuple(int(part) if part.isdigit() else 0 for part in parts)


def libmpv_target(environment):
    cache = Path(environment.get('PUB_CACHE') or Path.home() / '.pub-cache')
    plugins = list((cache / 'hosted' / 'pub.dev').glob('media_kit_libs_android_video-*'))
    if not plugins:
        return None
    return max(plugins, key=_plugin_version_key) / 'android' / 'build' / LIBMPV_VERSION


def prefetch_libmpv(environment):
    target = libmpv_target(environment)
    if target is None:
        print('未找到 media_kit_libs_android_video，跳过 libmpv 预下载。', flush=True)
        return
    target.mkdir(parents=True, exist_ok=True)
    for abi, digest in LIBMPV_JARS:
        destination = target / f'default-{abi}.jar'
        if destination.is_file() and hashlib.md5(destination.read_bytes()).hexdigest() == digest:
            continue
        destination.unlink(missing_ok=True)
        url = LIBMPV_URL.format(version=LIBMPV_VERSION, abi=abi)
        for attempt in range(1, 6):
            temporary = target / f'default-{abi}.part'
            try:
                print(f'下载 libmpv {abi}（第 {attempt} 次）：{url}', flush=True)
                request = urllib.request.Request(url, headers={'User-Agent': 'Mozilla/5.0'})
                with urllib.request.urlopen(request, timeout=180) as response:
                    temporary.write_bytes(response.read())
                if hashlib.md5(temporary.read_bytes()).hexdigest() != digest:
                    raise RuntimeError('MD5 校验失败')
                temporary.replace(destination)
                break
            except Exception as error:
                temporary.unlink(missing_ok=True)
                destination.unlink(missing_ok=True)
                if attempt == 5:
                    raise SystemExit(f'libmpv {abi} 下载失败，请检查网络或稍后重试：{error}')
                time.sleep(attempt * 5)

root = Path(__file__).resolve().parents[1]
parser = argparse.ArgumentParser()
parser.add_argument('--abi', action='append', choices=['arm64-v8a', 'armeabi-v7a', 'x86_64'])
parser.add_argument('--cn-mirrors', action='store_true', help='使用 Flutter 中国镜像和阿里云 Maven 镜像')
add_variant_argument(parser)
options = parser.parse_args()
variant = BuildVariant(options.all_sources)
environment = os.environ.copy()
if platform.system() == 'Darwin':
    environment['LANG'] = 'en_US.UTF-8'
    environment['LC_ALL'] = 'en_US.UTF-8'
environment.setdefault('GOPROXY', 'https://goproxy.cn,direct')
environment.setdefault('GOSUMDB', 'off')
flutter = shutil.which('flutter')
if not flutter:
    raise SystemExit('请先将 Flutter SDK 的 bin 目录加入 PATH。')
abi_args = [item for abi in options.abi or [] for item in ['--abi', abi]]
with china_mirror_environment(environment, options.cn_mirrors) as env, mirrored_pub_lockfile(root, env):
    if options.cn_mirrors:
        print('本次构建启用国内依赖镜像，保留官方 Maven 仓库备用。', flush=True)
    subprocess.run([sys.executable, str(root / 'scripts' / 'build_native.py'), '--platform', 'android', *abi_args, *variant.arguments],
                   cwd=root, env=env, check=True)
    subprocess.run([flutter, 'pub', 'get', '--enforce-lockfile'], cwd=root, env=env, check=True)
    prefetch_libmpv(env)
    build_args = [flutter, 'build', 'apk', '--release', '--split-per-abi', '--no-pub', *variant.flutter_arguments]
    if options.abi:
        targets = {'arm64-v8a': 'android-arm64', 'armeabi-v7a': 'android-arm', 'x86_64': 'android-x64'}
        build_args += ['--target-platform', ','.join(targets[abi] for abi in options.abi)]
    subprocess.run(build_args, cwd=root, env=env, check=True)
    subprocess.run([sys.executable, str(root / 'scripts' / 'package_release.py'), '--platform', 'android', *abi_args, *variant.arguments],
                   cwd=root, env=env, check=True)
