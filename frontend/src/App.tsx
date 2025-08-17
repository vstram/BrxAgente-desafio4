import {useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet, SelecionarDiretorio, SetDiretorioPlanilhas, RealizarAnaliseOrquestrada} from "../wailsjs/go/main/App";

function App() {
    const [resultText, setResultText] = useState("Por favor, selecione o diretório das planilhas abaixo 👇");
    const [name, setName] = useState('');
    const [diretorio, setDiretorio] = useState('');
    const [diretorioValido, setDiretorioValido] = useState(false);
    const [analiseEmAndamento, setAnaliseEmAndamento] = useState(false);
    
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

    async function fazerAnalise() {
        if (!diretorio || !diretorioValido) {
            setResultText("Por favor, selecione um diretório válido primeiro!");
            return;
        }

        setAnaliseEmAndamento(true);
        setResultText("Análise em andamento, por favor aguarde...");

        try {
            const resultado = await RealizarAnaliseOrquestrada(diretorio);
            setResultText(resultado);
        } catch (err: any) {
            setResultText(`Erro durante a análise: ${err.message}`);
        } finally {
            setAnaliseEmAndamento(false);
        }
    }

    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo"/>
            <div id="result" className="result">{resultText}</div>
            
            {/* Seção de seleção de diretório */}
            <div id="diretorio" className="input-box">
                <button className="btn" onClick={selecionarDiretorio}>Selecionar Diretório</button>
                {diretorio && (
                    <div>
                        <p>Diretório selecionado: {diretorio}</p>
                        <p>Validade: {diretorioValido ? 'Válido' : 'Inválido'}</p>
                    </div>
                )}
            </div>
            
            {/* Botão de análise - só aparece quando há um diretório válido selecionado */}
            {diretorioValido && (
                <div id="analise" className="input-box">
                    <button 
                        className="btn" 
                        onClick={fazerAnalise} 
                        disabled={analiseEmAndamento}
                    >
                        {analiseEmAndamento ? 'Analisando...' : 'Fazer Análise'}
                    </button>
                </div>
            )}
            
            {/* Seção de greeting - manter para testes, mas pode ser ocultada */}
            <div id="input" className="input-box" style={{display: 'none'}}>
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text"/>
                <button className="btn" onClick={greet}>Greet</button>
            </div>
        </div>
    )
}

export default App
