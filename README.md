# goapp
基于CEF3和Go语言的桌面应用程序框架。

# 如何编译 (Win7 32/64-bit)
1. Make sure you have installed windows-386 version of Go (for example: 1.3.3)
2. Install mingw and add C:\MinGW\bin to PATH. You can install mingw using mingw-get-setup.exe. Select packages to install: "mingw-developer-toolkit", "mingw32-base", "msys-base". CEF2go was tested and works fine with GCC 4.8.1. You can check gcc version with "gcc --version".
3. Download CEF 3 branch 2171 revision 1897 binaries: [cef_binary_3.2171.1897_windows32.7z](http://pan.baidu.com/s/1eQkYTYa)
   Copy Release/* to cef2go/Release
   Copy Resources/* to cef2go/Release
4. Run build.bat in the directory
5. Copy conf/manifest.json to bin directory
6. Run goapp.exe in bin directory

关于:
-------------------

作者 email: 529808348@qq.com




