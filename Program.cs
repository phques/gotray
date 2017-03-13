using Microsoft.Win32;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ConsoleApplication1
{
    class Program
    {
        static void Main(string[] args)
        {
            string keypath = @"Software\Classes\Local Settings\Software\Microsoft\Windows\CurrentVersion\TrayNotify";
            string valuename = "IconStreams";

            Console.WriteLine("load key HKCU\\" + keypath);

            var key = Registry.CurrentUser.OpenSubKey(keypath);
            if (key != null)
            {
                var valuekind = key.GetValueKind(valuename);
                if (valuekind == RegistryValueKind.Binary)
                {
                    byte[] value = key.GetValue(valuename) as byte[];
                }
                else
                {
                    Console.WriteLine("failed to load key" );
                }

                key.Close();
            }
            else
            {
                Console.WriteLine("key is not binary");
            }
        }
    }
}
