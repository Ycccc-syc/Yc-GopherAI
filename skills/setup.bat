@echo off
chcp 65001 >nul
REM ============================================
REM AI 全栈挑战赛 — 新项目初始化
REM 用法: .\skills\setup.bat D:\新项目路径
REM ============================================

if "%1"=="" (
    echo 用法: .\skills\setup.bat D:\目标项目路径
    exit /b 1
)

set TARGET=%1

if not exist "%TARGET%" mkdir "%TARGET%"

echo.
echo === 复制技能包到新项目 ===
echo.

xcopy /E /I /Y "%~dp0" "%TARGET%\skills\" >nul
echo [OK] skills/

copy /Y "%~dp0..\.claude\settings.json" "%TARGET%\.claude\settings.json" >nul 2>&1
echo [OK] .claude/settings.json  ^(权限配置^)

copy /Y "%~dp0..\CLAUDE.md" "%TARGET%\CLAUDE.md" >nul
echo [OK] CLAUDE.md  ^(记得修改内容适配新项目^)

echo.
echo === 完成 ===
echo.
echo Skills 在 VSCode 里全局有效，已不用再配。
echo 去新项目改好 CLAUDE.md 就可以开始了。
echo.
pause
