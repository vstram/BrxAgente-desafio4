import { useState, useRef, useEffect } from 'react';
import './Chat.css';
import { AskAI } from "../wailsjs/go/main/App";

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
  const messagesEndRef = useRef<HTMLDivElement>(null);

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
      // Get AI response
      const response = await AskAI(inputValue);
      
      // Add AI message
      const aiMessage: Message = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: response,
        timestamp: new Date(),
      };

      setMessages(prev => [...prev, aiMessage]);
    } catch (err: any) {
      // Add error message
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: `Desculpe, ocorreu um erro: ${err.message}`,
        timestamp: new Date(),
      };

      setMessages(prev => [...prev, errorMessage]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const clearChat = () => {
    setMessages([]);
  };

  return (
    <>
      {isOpen && (
        <div className="chat-overlay" onClick={() => setIsOpen(false)}>
          <div className="chat-window" onClick={(e) => e.stopPropagation()}>
            <div className="chat-header">
              <h3>Assistente de Dados</h3>
              <div className="chat-header-actions">
                <button className="chat-clear-btn" onClick={clearChat}>Limpar</button>
                <button className="chat-close-btn" onClick={() => setIsOpen(false)}>×</button>
              </div>
            </div>
            
            <div className="chat-messages">
              {messages.length === 0 ? (
                <div className="chat-welcome">
                  <p>Olá! Sou seu assistente de dados. Posso ajudar você a entender os resultados do processamento de VR/VA.</p>
                  <p>Exemplos de perguntas:</p>
                  <ul>
                    <li>Quantos colaboradores foram processados?</li>
                    <li>Qual o valor total de VR concedido?</li>
                    <li>Quais colaboradores estão com desligamento neste mês?</li>
                  </ul>
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
                onKeyPress={handleKeyPress}
                placeholder="Digite sua pergunta..."
                disabled={isLoading}
                rows={3}
              />
              <button 
                className="chat-send-btn" 
                onClick={handleSend} 
                disabled={isLoading || !inputValue.trim()}
              >
                Enviar
              </button>
            </div>
          </div>
        </div>
      )}
      
      <button className="chat-toggle-btn" onClick={() => setIsOpen(!isOpen)}>
        💬
      </button>
    </>
  );
}

export default Chat;