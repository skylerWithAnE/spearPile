spearePile.exe %1
echo %errorlevel% > Output
SET /p MYVAR=<Output
ECHO %MYVAR%
DEL Output