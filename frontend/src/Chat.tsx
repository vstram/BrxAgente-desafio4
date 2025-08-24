import { useState, useRef, useEffect } from 'react';
import './Chat.css';
import { AskAI, AskAIBasic, GetSystemPrompt, IsAgentEnabled, SetAgentEnabled } from "../wailsjs/go/main/App";

// Importando ícones
import ChatIcon from './assets/icons/chat.svg';
import SendIcon from './assets/icons/send.svg';
import ClearIcon from './assets/icons/clear.svg';
import XIcon from './assets/icons/x.svg';
import EyeIcon from './assets/icons/eye.svg';

interface Message {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: Date;
}

function Chat() {
  const [isOpen, setIsOpen] = useState(false);
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputValue, setInputValue] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isInspectorOpen, setIsInspectorOpen] = useState(false);
  const [systemPrompt, setSystemPrompt] = useState('');
  const [isAgentEnabled, setIsAgentEnabled] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Load initial agent state
  useEffect(() => {
    const loadAgentState = async () => {
      try {
        const enabled = await IsAgentEnabled();
        setIsAgentEnabled(enabled);
      } catch (err) {
        console.error('Error loading agent state:', err);
      }
    };
    loadAgentState();
  }, []);

  // Scroll to bottom of messages
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSend = async () => {
    if (!inputValue.trim() || isLoading) return;

    // Add user message
    const userMessage: Message = {
      id: Date.now().toString(),
      role: 'user',
      content: inputValue,
      timestamp: new Date(),
    };

    setMessages(prev => [...prev, userMessage]);
    setInputValue('');
    setIsLoading(true);

    try {
      // Get AI response - use agent if enabled, otherwise use basic mode
      const response = isAgentEnabled 
        ? await AskAI(inputValue)
        : await AskAIBasic(inputValue);
      
      // Add AI message
      const aiMessage: Message = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: response,
        timestamp: new Date(),
      };

      setMessages(prev => [...prev, aiMessage]);
    } catch (err: any) {
      console.error('Chat error:', err);
      
      // Add error message with better handling
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: `Desculpe, ocorreu um erro: ${err?.message || err?.toString() || 'Erro desconhecido'}`,
        timestamp: new Date(),
      };

      setMessages(prev => [...prev, errorMessage]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const clearChat = () => {
    setMessages([]);
  };

  const toggleAgent = async () => {
    try {
      const newState = !isAgentEnabled;
      await SetAgentEnabled(newState);
      setIsAgentEnabled(newState);
    } catch (err: any) {
      console.error('Error toggling agent:', err);
      // Show error message in chat
      const errorMessage: Message = {
        id: Date.now().toString(),
        role: 'assistant',
        content: `Erro ao alterar o estado do agente: ${err?.message || err?.toString() || 'Erro desconhecido'}`,
        timestamp: new Date(),
      };
      setMessages(prev => [...prev, errorMessage]);
    }
  };

  const inspectSystemPrompt = async () => {
    try {
      // Get the system prompt from the backend
      const prompt = await GetSystemPrompt();
      setSystemPrompt(prompt);
      setIsInspectorOpen(true);
    } catch (err: any) {
      setSystemPrompt(`Erro ao obter o System Prompt: ${err.message}`);
      setIsInspectorOpen(true);
    }
  };

  return (
    <>
      {isOpen && (
        <div className="chat-overlay" onClick={() => setIsOpen(false)}>
          <div className="chat-window" onClick={(e) => e.stopPropagation()}>
            <div className="chat-header">
              <h3>Assistente de Dados</h3>
              <div className="chat-header-actions">
                <button 
                  className={`chat-agent-toggle ${isAgentEnabled ? 'enabled' : 'disabled'}`} 
                  onClick={toggleAgent}
                  title={isAgentEnabled ? 'Agente IA Avançado Ativado' : 'Modo Básico Ativado'}
                >
                  {isAgentEnabled ? '🧠' : '💬'}
                </button>
                <button className="chat-inspect-btn" onClick={inspectSystemPrompt}>
                  <img src={EyeIcon} alt="Inspecionar" className="icon" />
                </button>
                <button className="chat-clear-btn" onClick={clearChat}>
                  <img src={ClearIcon} alt="Limpar chat" className="icon" />
                </button>
                <button className="chat-close-btn" onClick={() => setIsOpen(false)}>
                  <img src={XIcon} alt="Fechar" className="icon" />
                </button>
              </div>
            </div>
            
            <div className="chat-messages">
              {messages.length === 0 ? (
                <div className="chat-welcome">
                  <p>
                    Olá! Sou seu assistente de dados. 
                    {isAgentEnabled 
                      ? ' Estou funcionando em modo avançado com Agente IA 🧠' 
                      : ' Estou funcionando em modo básico 💬'
                    }
                  </p>
                  <p>Posso ajudar você a entender os resultados do processamento de VR/VA.</p>
                  
                  {isAgentEnabled ? (
                    // Modo Avançado - Funcionalidades IA
                    <>
                      <p><strong>🧠 Funcionalidades Avançadas do Agente IA (usando LangChainGO):</strong></p>
                      <p><small>💡 Use o botão {isAgentEnabled ? '🧠' : '💬'} no cabeçalho para alternar entre os modos Básico/Avançado.</small></p>
                      <ul>
                        <li><strong>Análise Inteligente de Padrões:</strong> Detecta automaticamente anomalias e inconsistências nos dados de VR/VA</li>
                        <li><strong>Respostas Contextuais:</strong> Compreende perguntas complexas e fornece análises detalhadas com base no histórico</li>
                        <li><strong>Recomendações Proativas:</strong> Sugere ações e melhorias baseadas na análise dos dados processados</li>
                      </ul>
                      <p>Exemplos de perguntas avançadas:</p>
                      <ul>
                        <li>Identifique anomalias nos valores de VR por sindicato</li>
                        <li>Compare a distribuição de dias úteis entre os sindicatos</li>
                        <li>Que recomendações você tem para otimizar os custos de VR?</li>
                      </ul>
                    </>
                  ) : (
                    // Modo Básico - Funcionalidades padrão
                    <>
                      <p>Exemplos de perguntas:</p>
                      <ul>
                        <li>Quantos colaboradores foram processados?</li>
                        <li>Qual o valor total de VR concedido?</li>
                        <li>Quais colaboradores estão com desligamento neste mês?</li>
                      </ul>
                    </>
                  )}
                  
                </div>
              ) : (
                messages.map((message) => (
                  <div key={message.id} className={`chat-message ${message.role}`}>
                    <div className="chat-message-content">
                      {message.content}
                    </div>
                    <div className="chat-message-time">
                      {message.timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                    </div>
                  </div>
                ))
              )}
              {isLoading && (
                <div className="chat-message assistant">
                  <div className="chat-message-content typing-indicator">
                    <span></span>
                    <span></span>
                    <span></span>
                  </div>
                </div>
              )}
              <div ref={messagesEndRef} />
            </div>
            
            <div className="chat-input-area">
              <textarea
                value={inputValue}
                onChange={(e) => setInputValue(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder="Digite sua pergunta..."
                disabled={isLoading}
                rows={3}
              />
              <button 
                className="chat-send-btn" 
                onClick={handleSend} 
                disabled={isLoading || !inputValue.trim()}
              >
                <img src={SendIcon} alt="Enviar" className="icon" />
              </button>
            </div>
          </div>
        </div>
      )}
      
      {isInspectorOpen && (
        <div className="inspector-overlay" onClick={() => setIsInspectorOpen(false)}>
          <div className="inspector-window" onClick={(e) => e.stopPropagation()}>
            <div className="inspector-header">
              <h3>System Prompt</h3>
              <button className="inspector-close-btn" onClick={() => setIsInspectorOpen(false)}>
                <img src={XIcon} alt="Fechar" className="icon" />
              </button>
            </div>
            <div className="inspector-content">
              <pre>{systemPrompt}</pre>
            </div>
          </div>
        </div>
      )}
      
      <button className="chat-toggle-btn" onClick={() => setIsOpen(!isOpen)}>
        <img src={ChatIcon} alt="Chat" className="icon" />
      </button>
    </>
  );
}

export default Chat;