package main

import (
	"bufio"
	"log"
	"fmt"
	"net/rpc"
	"os"
	"strings"
)

type Chat struct {
	Usuarios []string
	Mensajes [][]string
}

var nickname string

var ultimoMensaje int

func verificarMensajes(pClient *rpc.Client, reply Chat, pNick string) {

	for {

		pClient.Call("APP.ObtenerMensajes", "" , &reply)

		msg := reply.Mensajes

		for i := ultimoMensaje; i <= len(msg)-1; i++ {

			if i == ultimoMensaje {
				continue
			} else if msg[i][1] == pNick && msg[i][2] == "0" {
				log.Println("Tu: " +  msg[i][0])
			} else if msg[i][1] != pNick && msg[i][2] == "0" {
				log.Println(msg[i][1] +  " dice: " +  msg[i][0])
			} else {
				log.Println(msg[i][0])
			}

		}

		ultimoMensaje = len(msg)-1

	}

}

func mainLoop( pLector *bufio.Reader, pClient *rpc.Client, pNick string, pReply Chat ) {
	
	reply := pReply
	nick := pNick
	lector :=  pLector

	ultimoMensaje = 0

	listarUsuarios(pClient, reply)

	go verificarMensajes(pClient, reply, nick )

	for {

		entrada, error := lector.ReadString('\n')
		entrada = strings.TrimSpace(entrada)

		if error != nil {
			log.Printf("Error: %q\n", error)
		}

		if strings.TrimSpace(entrada) == "" {

		
		} else if strings.HasPrefix(entrada, "/help") || strings.HasPrefix(entrada, "/HELP") {

			fmt.Println("")
			fmt.Println("'/salir' - Abandonas el chat.")
			fmt.Println("'/usuarios' - Lista los usuarios que estan actualmente ON.")
			fmt.Println("'-texto-' - Envias un mensaje al chat.")
			fmt.Println("")
			
		} else if strings.HasPrefix(entrada, "/salir") || strings.HasPrefix(entrada, "/SALIR") {

			fmt.Println("Has abandonado el chat.")
			pClient.Call("APP.UsuarioSalir", nick , &reply)
			break
		} else if strings.HasPrefix(entrada, "/usuarios") || strings.HasPrefix(entrada, "/usuarios") {

			listarUsuarios(pClient, reply)

		} else {

			pClient.Call("APP.RegistrarMensaje", []string {entrada, nick, "0"} , &reply)
		}
	}	
}

func listarUsuarios(pClient *rpc.Client, pReply Chat) {

	reply := pReply

	fmt.Println("")

	pClient.Call("APP.ObtenerDatos", "" , &reply)

	fmt.Println("Usuarios ON:")
	for i := range reply.Usuarios {

		if reply.Usuarios[i] == nickname {
			fmt.Println(reply.Usuarios[i] + " (Yo)") 
		} else {
			fmt.Println(reply.Usuarios[i])
		} 
	}

	fmt.Println("")

}


func main() {

	var reply Chat

	client, err := rpc.DialHTTP("tcp", "localhost:4040")

	if err != nil {
		log.Fatal("Error de conexion: ", err)
	}

	lector := bufio.NewReader(os.Stdin)

	log.Println("Ingrese su NickName:")
	nickname, err :=  lector.ReadString('\n')
	nickname = strings.TrimSpace(nickname)

	if err != nil {
		log.Fatal("Error con el username: ", err)
	}

	err = client.Call("APP.UsuarioExiste", nickname , &reply )

	if err != nil {
		log.Println("ERROR")
		log.Fatal( err )
	} 

	client.Call("APP.RegistrarUsuario", nickname , &reply )
	log.Println("Bienvenido " + nickname + "\n")

	fmt.Println("")

	fmt.Println("Si quieres salir del chat, ingresa el comando '/salir'")
	fmt.Println("Para ver los comandos existentes ingresa '/help'")

	mainLoop(lector, client, nickname, reply)


}
