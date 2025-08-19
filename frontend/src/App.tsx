import {useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet, SelecionarDiretorio, SetDiretorioPlanilhas, RealizarAnaliseOrquestrada} from "../wailsjs/go/main/App";
import ConfigModal from './ConfigModal';
import Chat from './Chat';

function App() {
    const [resultText, setResultText] = useState("Por favor, selecione o diretório das planilhas abaixo 👇");
    const [name, setName] = useState('');
    const [diretorio, setDiretorio] = useState('');
    const [diretorioValido, setDiretorioValido] = useState(false);
    const [analiseEmAndamento, setAnaliseEmAndamento] = useState(false);
    const [resultados, setResultados] = useState<{colaboradoresProcessados?: number, arquivoGerado?: string} | null>(null);
    const [erro, setErro] = useState<string | null>(null);
    const [isConfigModalOpen, setIsConfigModalOpen] = useState(false);
    
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
                setResultados(null);
                setErro(null);
                
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
        setResultados(null);
        setErro(null);

        try {
            const resultado = await RealizarAnaliseOrquestrada(diretorio);
            setResultText(resultado);
            
            // Extrair informações do resultado para exibição formatada
            const match = resultado.match(/(\d+) colaboradores processados.*como (.+\.xlsx)/);
            if (match) {
                setResultados({
                    colaboradoresProcessados: parseInt(match[1]),
                    arquivoGerado: match[2]
                });
            }
        } catch (err: any) {
            const errorMessage = `Erro durante a análise: ${err.message}`;
            setResultText(errorMessage);
            setErro(errorMessage);
        } finally {
            setAnaliseEmAndamento(false);
        }
    }

    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo"/>
            <div id="result" className="result">{resultText}</div>
            
            {/* Botão de configuração */}
            <div id="config" className="input-box">
                <button className="btn" onClick={() => setIsConfigModalOpen(true)}>Configurações</button>
            </div>
            
            {/* Área de exibição de resultados */}
            <div id="resultados" className="resultados-box">
                {analiseEmAndamento && (
                    <div className="status-em-andamento">
                        <p>Processando dados... ⏳</p>
                    </div>
                )}
                
                {erro && (
                    <div className="status-erro">
                        <h3>Erro no Processamento</h3>
                        <p>{erro}</p>
                    </div>
                )}
                
                {resultados && (
                    <div className="status-concluido">
                        <h3>Análise Concluída com Sucesso!</h3>
                        <p><strong>Colaboradores processados:</strong> {resultados.colaboradoresProcessados}</p>
                        <p><strong>Arquivo gerado:</strong> {resultados.arquivoGerado}</p>
                        <p>O arquivo foi salvo na pasta Downloads.</p>
                    </div>
                )}
            </div>
            
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
            
            {/* Modal de configuração */}
            <ConfigModal 
                isOpen={isConfigModalOpen} 
                onClose={() => setIsConfigModalOpen(false)} 
            />
            
            {/* Componente de chat */}
            <Chat />
        </div>
    )
}

export default App
