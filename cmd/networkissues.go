package cmd

import (
  "golang.org/x/sys/windows/registry"
  "syscall"
)

func checkProxy(isVpn bool) (bool, string, string, error) {   /* isProxy, proxy, problem, err */
  //isVpn, err := cmd.Flags().GetBool("vpn")
  var proxyStr = ""

  k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.QUERY_VALUE)
  if err != nil {
    //log.Fatal(err)
    return false, proxyStr, "APIFAIL", err
  }
  defer k.Close()

  // -- auto config URL -----------
  autoConfigUrl, _, errAutoConfig := k.GetStringValue("AutoConfigURL")
  if errAutoConfig != nil && errAutoConfig != syscall.ENOENT {
    //log.Fatal(errAutoConfig)
    return false, proxyStr, "APIFAIL", err
  }

  // -- proxy enabled -----------
  proxyEnabled, _, err := k.GetIntegerValue("ProxyEnable")
  if err != nil {
    //log.Fatal(err)
    return false, proxyStr, "APIFAIL", err
  }

  // -- proxy server -----------
  proxyServer, _, err := k.GetStringValue("ProxyServer")
  if err != nil {
    //log.Fatal(err)
    return false, proxyStr, "APIFAIL", err
  }

  // What is the proxy?
  if proxyEnabled == 1 && len(proxyServer) > 0 {
    proxyStr = autoConfigUrl
  }
  if errAutoConfig != syscall.ENOENT && len(autoConfigUrl) != 0 {
    proxyStr = autoConfigUrl
  }


  // -- Is there a proxy issue -----------
  numProxies := 0
  if isVpn {
    if errAutoConfig != syscall.ENOENT && len(autoConfigUrl) != 0 {
      numProxies += 1
    }

    if proxyEnabled == 1 && len(proxyServer) > 0 {
      numProxies += 1
    }

    if numProxies == 0 {
      //fmt.Println("Bad -- On VPN, but no autoConfigUrl or proxy")
      return false, proxyStr, "Bad -- On VPN, but no autoConfigUrl or proxy", nil
    } else if numProxies > 1 {
      //fmt.Println("Bad -- On VPN, but both autoconfig and proxy")
      return true, proxyStr, "Bad -- On VPN, but both autoconfig and proxy", nil
    }

  } else {
    if len(autoConfigUrl) > 0 {
      // This is bad, generally.
      //fmt.Println("Bad -- no-VPN but have autoConfigUrl:", autoConfigUrl)
      return true, proxyStr, "Bad -- no-VPN but have autoConfigUrl: " + autoConfigUrl, nil
    }

    if proxyEnabled == 1 && len(proxyServer) > 0 {
      //fmt.Println("Bad -- no-VPN but have proxy:", proxyServer)
      return true, proxyStr, "Bad -- no-VPN but have proxy:" + proxyServer, nil
    }
  }

  return numProxies > 0, proxyStr, "", nil
}

