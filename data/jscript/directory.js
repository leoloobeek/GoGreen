function WalkOS(allPaths,fso,folder,depth) {
    var folEnum = new Enumerator(folder.SubFolders);
    for (;!folEnum.atEnd(); folEnum.moveNext()) {
        allPaths.push(fso.GetAbsolutePathName(folEnum.item()));
        if(depth > 0) {
            WalkOS(allPaths,fso,folEnum.item(),(depth-1));
        }
    }
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
    WalkOS(allPaths,objFSO,startDirObject,~DEPTH~);  
}
catch(err) {}