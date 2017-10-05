Function WalkOS(allPaths,fso,folder)
    dim files
    allPaths.add(fso.GetAbsolutePathName(folder))    
    For Each Subfolder In folder.SubFolders
        On Error Resume Next
        WalkOS allPaths,fso,Subfolder
        If Err.Number <> 0 Then
            return
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
WalkOS allPaths,objFSO,startDirObject