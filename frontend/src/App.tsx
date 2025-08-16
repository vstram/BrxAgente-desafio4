import {useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet, SelecionarDiretorio, SetDiretorioPlanilhas} from "../wailsjs/go/main/App";

function App() {
    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [name, setName] = useState('');
    const [diretorio, setDiretorio] = useState('');
    const [diretorioValido, setDiretorioValido] = useState(false);
    
    const updateName = (e: any) => setName(e.target.value);
    const updateResultText = (result: string) => setResultText(result);

    function greet() {
        Greet(name).then(updateResultText);
    }

    async function selecionarDiretorio() {
        try {
            const caminho = await SelecionarDiretorio();
            if (caminho) {
                setDiretorio(caminho);
                
                // Validando o diretório no backend
                try {
                    const valido = await SetDiretorioPlanilhas(caminho);
                    setDiretorioValido(valido);
                    if (valido) {
                        setResultText(`Diretório válido selecionado: ${caminho}`);
                    } else {
                        setResultText(`Diretório inválido: ${caminho}`);
                    }
                } catch (err: any) {
                    setDiretorioValido(false);
                    setResultText(`Erro ao validar diretório: ${err.message}`);
                }
            }
        } catch (err: any) {
            setResultText(`Erro ao selecionar diretório: ${err.message}`);
        }
    }

    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo"/>
            <div id="result" className="result">{resultText}</div>
            <div id="input" className="input-box">
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text"/>
                <button className="btn" onClick={greet}>Greet</button>
            </div>
            <div id="diretorio" className="input-box">
                <button className="btn" onClick={selecionarDiretorio}>Selecionar Diretório</button>
                {diretorio && (
                    <div>
                        <p>Diretório selecionado: {diretorio}</p>
                        <p>Validade: {diretorioValido ? 'Válido' : 'Inválido'}</p>
                    </div>
                )}
            </div>
        </div>
    )
}

export default App
