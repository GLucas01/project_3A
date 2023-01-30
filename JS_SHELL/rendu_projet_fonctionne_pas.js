import inq from 'inquirer';
import boxen from 'boxen';
import psList from 'ps-list';
import cp from 'child_process';
import { exec } from 'child_process'
import * as process from "process";
import { suspend, resume } from 'ntsuspend';
import chalk from 'chalk';
import fs from 'fs';


// variables
var cmdhistorytab = [];
var running = true;
var lastcommand = null;
var mainpath = process.cwd();  // connaître le répertoire de travail actuel
const actions = [
    {//une fonction dans l'array pour avoir des types
        name: 'clear',
        desc: 'clears the console'
    },
    { 
        name: 'lp',
        desc: 'lists the running processes on the machine'
    },
    {
        name: 'exit',
        desc: 'exit the CLI. Ctrl+P also works'
    },
    {
        name: 'cd (directory)',
        desc: 'navigate through directories'
    },
    {
        name: 'exec  (path-to-program)',
        desc: 'runs a program from PATH variables or direct path'
    },
    {
        name: 'bing (-k|-p|-c) (process_id)',
        desc: "-k kills select process -p pauses and -c resumes"
    },
    ]

// commencement 
console.clear();
 
function saisie() {      // permettre à l'user de saisir une commande
    return inq.prompt(      // line utilise prompt de la bibli inq pour afficher une boîte de dialogue
        [
            {
                name: 'command',
                type: 'input',
                message: (mainpath + ' $ ')   // affiche  le répertoire actuel et $ indique que c'est une invite de commande

            }
        ]);
}

//fonctionenment
const run = async () => {  //fction asynchrone run pour traiter une entrée utilisateur

    exitCommand();     // vérifier si exit a été appelé 
    mainpath = process.cwd();   // maj de mainpath
    const cmd = await saisie();    //saisie attends une entrée qu'il affecte à cmd
     await action(cmd);   // action() execute la commande entrée
}


//actions
async function action (cmd) {      // fction des différentes actions en fonction de la commande.
    if(cmd.command === "clear"){
        console.clear();
    }

    else if(/^exec /.test(cmd.command)) { // si la commande est exec
        let prog = cmd.command.replace(/^exec /, "");   // enlève exec pour garder le nom du programme
        

        // Execute la commande
        await cp.exec(prog, (error) => {         //utilise exec de child process pour executer le programme demandé
            if (error) {
                console.error(`exec error: ${error}`);  // affiche une potentielle erreur lors de l'execution
            }
        });
    }

    else if(cmd.command === "lp"){

        console.log((boxen('processus en cours ')));

      let processes = (await psList());
      processes.sort(compare);  // trier les processus par ordre croissant

      console.log('\n');
      for(let i = 0; i < 50; i++){  //affiche les 50 premiers
          console.log(processes[i].name+' pid: ' + chalk.red(processes[i].pid) +"\n");   //affiche les processus avec leur nom et leur pid
      }
    }

    

    else if(/^bing/.test(cmd.command)) {           // si la commande commence par bing 

        const com = cmd.command.match(/(-k|-p|-c) /);           // utilise match pour extraire k ou p ou c          
        const processId = cmd.command.replace(/^bing (-k|-p|-c) /,"");     // récupère le processId

            if(com != null) {
                switch (com[0]) {         //vérifie l'option de commande
                    case '-k ':  
                        // Kill le processus

                        try{
                            process.kill(processId);
                            console.log("Processus killed");
                        }catch{
                            console.log("Processus inexistant")
                        }
                        break;
                    case '-p ':
                        // Pause le processus
                        console.log("Processus paused: " );
                        if (process.platform === 'win32') {  //si sous windows
                            suspend(processId);
                        }else{
                            try {
                                process.kill(processId, 'SIGSTOP');           // sinon envois signal SIGSTOP 
                            }catch{
                                console.log("Processus inexistant.")
                            }
                        }
                        break;
                    case '-c ':
                        // Resume the process
                        console.log("Processus resumed: ");
                        if (process.platform === 'win32') {     //si sous windows
                            resume(processId);
                        }else{
                            process.kill(processId, 'SIGCONT');    //sinon envois signal
                        }
                        break;

                }
            } else {
                console.error('réessayez avec bing -k|-p|-c (process_id)');
            }
        }

    else if(/^cd(..$| )/.test(cmd.command)) {   // vérifie si l'entrée est cd .. ou cd 

        let pathname = cmd.command.replace(/^cd ?/,"");    // enlève cd et récupère le nom du path

        if(process.platform !== 'win32'){        // vérifie si on est sous windows
            if(pathname === '..'){
                pathname = '../';       // remonter d'un niveau dans le répertoire
            }
        }

        try{
            process.chdir(pathname);       // si on a pas .. mais le nom du répertoire alors on utilise chdir
        }catch (e) {
            console.log("Please enter a valid path !");   // si on fait une erreur de nom de chemin
        }
    }


   

    else if (cmd.command === 'exit') {
        process.exit();           // sors du processus
        running = false;          // mets à jour running
    }


    else {console.log("Unrecognized command:" + cmd.command)}    // cas terminal 
}

function exitCommand(){                  // sortie avec ctrl+p
    var stdin = process.stdin;

    stdin.setRawMode(true);   // passer en mode brut: lire les entrées clavier comme le système d'exploitation les envoi 
    stdin.resume();     // redémarre l'écoute des entrées clavier
    stdin.setEncoding('utf8');  //définir encodage utilisé pour les entrées clavier : encodage pour caractères.

    stdin.on('data', key => {
        // Ctrl+P
        if (key === '\u0010') {
            process.exit();   //quitter le programme si la commande est ctrl+p
            running = false;  
            console.log("bug");
        }
    });
}




