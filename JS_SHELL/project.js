let stop_boucle = false;

console.log("Welcome to ZK & Lucas Shell")
console.log("Type <<--help>> to discover the program")


while (!stop_boucle) {

    const prompt = require('prompt-sync')();
    const execSync = require('child_process').execSync;
    let command = prompt('+ :');
      
    if (command=="--help"){
        console.log("\n-------------------------------------")
        console.log("Type <<lp>> to show active processor")
        console.log("Type <<bing man>> to show bing manual")
        console.log("Type <<shutdown>> to quit the program")
        console.log("Type <<clear>> to clear the terminal")
        console.log("Type <<exec [program_name]>> to run a program in the operating system")
        console.log("-------------------------------------\n")

    }

    else if (command=="shutdown"){
        stop_boucle = true
     }

    else if (command=="lp"){
        const output = execSync('ps -a', { encoding: 'utf-8' });
        console.log(output);
    }
    else if (command.match('exec ([^\']+)')){
        program=command.match('exec ([^\']+)');
        console.log("Executing the program : "+program[1])
        execSync('open '+program[1], { encoding: 'utf-8' });

    }

    else if(command == "clear"){
        console.clear();
    }
    else if (command.match('bing ([^\']+)')){
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
            console.log("Invalid command. Try <<man bing>> for bing manual")
        }
    }

    else {
        console.log("Invalid command. Try <<--help>> for user manual")
    }

    
}
