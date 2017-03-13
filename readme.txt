article https://tmintner.wordpress.com/2011/07/08/windows-7-notification-area-automation-falling-back-down-the-binary-registry-rabbit-hole/

how to change the visibility flag of app icon in the Windows notification tray

original script : setvisible.ps1

The notification settings for the task tray are stored in the registry at 
HKCU\Software\Classes\Local Settings\Microsoft\Windows\CurrentVersion\TrayNotify 
in the IconStreams value as a binary registry key. 
Luckly for us, the organization of the key is not nearly as hard to understand as the Favorites Bar.
The binary stream begins with a 20 byte header followed by X number of 1640 byte items where X 
is the number of items that have notification settings.
Each 1640 byte block is comprised of at least
(one of the sections is not fully decoded so it may be made up of 2 or more sections)
5 fixed byte width sections as follows:

528 bytes – Path to the executable
4 bytes – Notification visibility setting
512 bytes – Last visible tooltip
592 bytes – Unknown (Seems to have a second tool-tip embeded in it but the starting position in the block changes)
4 bytes – ID?

0 = only show notifications
1 = hide icon and notifications
2 = show icon and notifications

*** Note that the values are actually saved by explorer.exe on logout and read on login,
so for the changes to be effective, one needs to kill explorer and re-login