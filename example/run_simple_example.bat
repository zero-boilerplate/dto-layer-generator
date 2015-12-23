@echo off
cls
SET ERRORLEVEL=0

::——————————————————— This will make it quit immediately upon pressing Ctrl+C instead of first asking yes/no
IF "%~1"=="–FIX_CTRL_C" (
SHIFT
) ELSE (
CALL <NUL %0 –FIX_CTRL_C %*
GOTO EOF
)
::———————————————————

pushd ".."
echo Running go install... ^
    & go install^
    & if errorlevel 1 goto ERROR
echo Running exe dto-layer-generator.exe from gopath/bin... ^
    & "%GOPATH%/bin/dto-layer-generator.exe" "example/simple_example.yml" ^
    & if errorlevel 1 goto ERROR

goto SUCCESS

:SUCCESS
popd
echo Success!!
goto EOF

:ERROR
popd
echo ERROR!!! See the last ran command
pause
goto EOF

:EOF