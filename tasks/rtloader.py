"""
RtLoader namespaced tasks
"""
import os

from invoke import task

def get_rtloader_path():
    here = os.path.abspath(os.path.dirname(__file__))
    return os.path.join(here, '..', 'rtloader')

def clear_cmake_cache(rtloader_path, settings):
    """
    CMake is not regenerated when we change an option. This function detect the
    current cmake settings and remove the cache if they have change to retrigger
    a cmake build.
    """
    cmake_cache = os.path.join(rtloader_path, "CMakeCache.txt")
    if not os.path.exists(cmake_cache):
        return

    settings = settings.copy()
    with open(cmake_cache) as cache:
        for line in cache.readlines():
            for key, value in settings.items():
                if line.strip() == key + "=" + value:
                    settings.pop(key)

    if settings:
        os.remove(cmake_cache)

@task
def build(ctx, install_prefix=None, python_runtimes=None, cmake_options=''):
    rtloader_path = get_rtloader_path()

    here = os.path.abspath(os.path.dirname(__file__))
    dev_path = os.path.join(here, '..', 'dev')

    cmake_args = cmake_options + " -DBUILD_DEMO:BOOL=OFF -DCMAKE_INSTALL_PREFIX:PATH={}".format(install_prefix or dev_path)

    python_runtimes = python_runtimes or os.environ.get("PYTHON_RUNTIMES") or "2"
    python_runtimes = python_runtimes.split(',')

    settings = {
            "DISABLE_PYTHON2:BOOL": "OFF",
            "DISABLE_PYTHON3:BOOL": "OFF"
            }
    if '2' not in python_runtimes:
        settings["DISABLE_PYTHON2:BOOL"] = "ON"
    if '3' not in python_runtimes:
        settings["DISABLE_PYTHON3:BOOL"] = "ON"

    # clear cmake cache if settings have changed since the last build
    clear_cmake_cache(rtloader_path, settings)

    for option, value in settings.items():
        cmake_args += " -D{}={} ".format(option, value)

    ctx.run("cd {} && cmake {} .".format(rtloader_path, cmake_args))
    ctx.run("make -C {}".format(rtloader_path))

@task
def install(ctx):
    ctx.run("make -C {} install".format(get_rtloader_path()))

@task
def test(ctx):
    ctx.run("make -C {}/test run".format(get_rtloader_path()))

@task
def format(ctx):
    ctx.run("make -C {} clang-format".format(get_rtloader_path()))
