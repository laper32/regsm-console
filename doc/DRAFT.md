```cs
// Server force stop previous existing servers.

// hlds.exe -game cstrike -console -ip 0.0.0.0 -port 23333 +map de_dust2 -maxplayer 32
var rootpath = @"C:\Users\user\Downloads\cs16";
var _path = Path.Combine(rootpath, "hlds.exe");

var exename = Path.GetFileNameWithoutExtension(_path);
var running = 0u;
Process.GetProcessesByName(exename).ToList().ForEach(p =>
{
    if (p.MainModule.FileName.Equals(_path))
    {
        if (!OperatingSystem.IsWindows()) return;

        Win32API.Message.Write(p.MainWindowHandle, "exit");
        Win32API.Message.Send(p.MainWindowHandle);
        
        p.WaitForExit(TimeSpan.FromSeconds(5).Milliseconds);

        if (!p.HasExited)
        {
            p.Kill();
        }
        p.WaitForExit();
        running++;
    }
});
Console.WriteLine($"We have closed {running} existing processes...");

Global.Application.ExecutableProcess = new()
{
    StartInfo =
    {
        FileName = _path,
        WorkingDirectory = rootpath,
        Arguments = "-game cstrike -console -ip 0.0.0.0 -port 27015 +map de_dust2 -maxplayer 32",
        WindowStyle = ProcessWindowStyle.Normal
    },
    EnableRaisingEvents = true,
};
Global.Application.ExecutableProcess.Start();
Console.ReadLine();


```