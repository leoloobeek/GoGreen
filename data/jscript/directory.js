function WalkOS(allPaths,fso,folder) {
    var folEnum = new Enumerator(folder.SubFolders);
    for (;!folEnum.atEnd(); folEnum.moveNext()) {
        WalkOS(allPaths,fso,folEnum.item());
    }
    allPaths.push(fso.GetAbsolutePathName(folder));
    var files = folder.Files;
    var fc = new Enumerator(folder.Files);
    for (;!fc.atEnd(); fc.moveNext()) {
        allPaths.push(fso.GetAbsolutePathName(fc.item()));
    }
}
var startDir = "~STARTDIR~";
try {
    var objFSO = new ActiveXObject("Scripting.FileSystemObject");
    var startDirObject = objFSO.GetFolder(startDir);
    WalkOS(allPaths,objFSO,startDirObject);  
}
catch(err) {}