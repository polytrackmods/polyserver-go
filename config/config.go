package config

var PolyVersion = "0.6.0-beta1"
var ApiVersion = "v6"

var IceFetchUrl = "https://vps.kodub.com:43274/" + ApiVersion + "/iceServers?version=" + PolyVersion
var WebsocketUrl = "wss://vps.kodub.com:43274/" + ApiVersion + "/multiplayer/host"

var AcceptVanillaClients = true

var LoadedMods []string = []string{}
