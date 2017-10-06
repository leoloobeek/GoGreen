function Get-EnvCombos {
        param($chars)
        $script:result = @()

        function Get-Combos
        {
            param($p,$c)
            if ($c.Length -eq 0) { break }
            For ($i=0; $i -le $c.Length; $i++) {
                $script:result += $p + $c[$i]
                Get-Combos "$p$($c[$i])" ($c[($i+1)..$c.Length])
            }
        }
        Get-Combos '' $chars -PassThru
        return $script:result
    }

    $oEnv = [Environment]::GetEnvironmentVariable
    $envCombos += Get-EnvCombos @(~ENVVARS~)