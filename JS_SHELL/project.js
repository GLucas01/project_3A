let boucle = false;


while (!boucle) {
	//const prompt = require('prompt-sync')({sigint: true});
	//sigint = signal interupt
	const prompt = require('prompt-sync')();
    const execSync = require('child_process').execSync;
    
    let command = prompt('+ :');
      
    if (command=="shutdown"){
        boucle = true
        }
    if (command=="lp"){
        const output = execSync('ps -a', { encoding: 'utf-8' });
        console.log(output);
    }
    if (command.match('bing ([^\']+)')){
        action=command.match('bing ([^\']+)');
        if (action[1].match('-k ([^\']+)')){
            console.log("Kill %d",action[1].match('-k ([^\']+)')[1])
            execSync('kill -9 '+String(action[1].match('-k ([^\']+)')[1]), { encoding: 'utf-8' });
        }
        else if (action[1].match('-p ([^\']+)')){
            console.log("Pause %d",action[1].match('-p ([^\']+)')[1])
            execSync('kill -17 '+String(action[1].match('-p ([^\']+)')[1]), { encoding: 'utf-8' });
        }
        else if (action[1].match('-c ([^\']+)')){
            console.log("Continue %d",action[1].match('-c ([^\']+)')[1])
            execSync('kill -s SIGCONT '+String(action[1].match('-c ([^\']+)')[1]), { encoding: 'utf-8' });
        }
        else if (action[1].match('man')){
            console.log("bing [-k|-p|-c] <processId>")
            console.log("-k for kill | -p for pause | -c for continue")
        }
        else {
            console.log("Invalid command. Try <<man bing>> for user manual")
        }
    }

    
}



