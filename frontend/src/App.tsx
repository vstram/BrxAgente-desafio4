import {useState} from 'react';
import logo from './assets/images/logo-vr-va.png';
import './App.css';
import {Greet, SelecionarDiretorio, SetDiretorioPlanilhas, RealizarAnaliseOrquestrada} from "../wailsjs/go/main/App";
import ConfigModal from './ConfigModal';
import Chat from './Chat';

// Importando ícones
import FolderIcon from './assets/icons/folder.svg';
import PlayIcon from './assets/icons/play.svg';
import SettingsIcon from './assets/icons/settings.svg';
import SpinnerIcon from './assets/icons/spinner.svg';

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
            {/* Cabeçalho */}
            <header className="app-header">
                <img src={logo} id="logo" alt="logo" />
                <h1>Automação de VR/VA</h1>
                <p className="app-subtitle">Sistema de processamento automatizado de Vale Refeição e Vale Alimentação</p>
            </header>
            
            {/* Área principal */}
            <main className="app-main">
                {/* Área de seleção de diretório */}
                <section className="section directory-section">
                    <h2>Seleção de Planilhas</h2>
                    <div className="section-content">
                        <button className="btn primary-btn" onClick={selecionarDiretorio}>
                            <img src={FolderIcon} alt="Selecionar diretório" className="btn-icon" />
                            Selecionar Diretório
                        </button>
                        {diretorio && (
                            <div className="directory-info">
                                <p><strong>Diretório selecionado:</strong> {diretorio}</p>
                                <p className={`status ${diretorioValido ? 'valid' : 'invalid'}`}>
                                    <strong>Validade:</strong> {diretorioValido ? 'Válido' : 'Inválido'}
                                </p>
                            </div>
                        )}
                    </div>
                </section>
                
                {/* Botão de análise */}
                {diretorioValido && (
                    <section className="section analysis-section">
                        <h2>Processamento</h2>
                        <div className="section-content">
                            <button 
                                className="btn secondary-btn" 
                                onClick={fazerAnalise} 
                                disabled={analiseEmAndamento}
                            >
                                {analiseEmAndamento ? (
                                    <>
                                        <img src={SpinnerIcon} alt="Processando" className="btn-icon spinner" />
                                        Processando...
                                    </>
                                ) : (
                                    <>
                                        <img src={PlayIcon} alt="Iniciar processamento" className="btn-icon" />
                                        Iniciar Processamento
                                    </>
                                )}
                            </button>
                        </div>
                    </section>
                )}
                
                {/* Área de resultados */}
                <section className="section results-section">
                    <h2>Resultados</h2>
                    <div className="section-content">
                        <div id="result" className="result">{resultText}</div>
                        
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
                    </div>
                </section>
            </main>
            
            {/* Rodapé */}
            <footer className="app-footer">
                <button className="btn footer-btn" onClick={() => setIsConfigModalOpen(true)}>
                    <img src={SettingsIcon} alt="Configurações" className="btn-icon" />
                    Configurações
                </button>
            </footer>
            
            {/* Seção de greeting - manter para testes, mas pode ser ocultada */}
            <div id="input" className="input-box" style={{display: 'none'}}>
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text" />
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
