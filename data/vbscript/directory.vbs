Function WalkOS(allPaths,fso,folder,depth)
    dim files
    For Each Subfolder In folder.SubFolders
        allPaths.add(fso.GetAbsolutePathName(Subfolder))
        If depth > 0 Then
            On Error Resume Next
            WalkOS allPaths,fso,Subfolder,(depth-1)
            If Err.Number <> 0 Then
                return
            End If
        End If
    Next
    For Each file In folder.Files
        allPaths.add(fso.GetAbsolutePathName(file))
    Next
End Function

dim startDir, allPaths, objFSO, startDirObject
startDir = "~STARTDIR~"
Set objFSO = CreateObject("Scripting.FileSystemObject")

Set startDirObject = objFSO.GetFolder(startDir)
WalkOS allPaths,objFSO,startDirObject,~DEPTH~