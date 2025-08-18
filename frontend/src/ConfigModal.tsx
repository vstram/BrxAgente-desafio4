import { useState, useEffect } from 'react';
import './ConfigModal.css';
import { GetConfig, SetOpenAIKey, SetOllamaConfig, TestOpenAIKey, TestOllamaConnection } from "../wailsjs/go/main/App";

interface ConfigModalProps {
  isOpen: boolean;
  onClose: () => void;
}

function ConfigModal({ isOpen, onClose }: ConfigModalProps) {
  const [openAIKey, setOpenAIKey] = useState('');
  const [ollamaURL, setOllamaURL] = useState('');
  const [ollamaModel, setOllamaModel] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [message, setMessage] = useState<{type: 'success' | 'error' | 'info', text: string} | null>(null);

  // Load current configuration when modal opens
  useEffect(() => {
    if (isOpen) {
      loadConfig();
    }
  }, [isOpen]);

  const loadConfig = async () => {
    try {
      setIsLoading(true);
      const config = await GetConfig();
      setOpenAIKey(config.openai_key || '');
      setOllamaURL(config.ollama_config?.base_url || '');
      setOllamaModel(config.ollama_config?.model || '');
    } catch (err: any) {
      setMessage({type: 'error', text: `Erro ao carregar configurações: ${err.message}`});
    } finally {
      setIsLoading(false);
    }
  };

  const handleSave = async () => {
    try {
      setIsLoading(true);
      setMessage(null);
      
      // Save OpenAI key
      await SetOpenAIKey(openAIKey);
      
      // Save Ollama config
      await SetOllamaConfig({
        base_url: ollamaURL,
        model: ollamaModel
      });
      
      setMessage({type: 'success', text: 'Configurações salvas com sucesso!'});
    } catch (err: any) {
      setMessage({type: 'error', text: `Erro ao salvar configurações: ${err.message}`});
    } finally {
      setIsLoading(false);
    }
  };

  const handleTestOpenAI = async () => {
    try {
      setIsLoading(true);
      setMessage(null);
      
      const result = await TestOpenAIKey(openAIKey);
      if (result) {
        setMessage({type: 'success', text: 'Chave da API do OpenAI é válida!'});
      } else {
        setMessage({type: 'error', text: 'Chave da API do OpenAI é inválida!'});
      }
    } catch (err: any) {
      setMessage({type: 'error', text: `Erro ao testar chave: ${err.message}`});
    } finally {
      setIsLoading(false);
    }
  };

  const handleTestOllama = async () => {
    try {
      setIsLoading(true);
      setMessage(null);
      
      const result = await TestOllamaConnection({
        base_url: ollamaURL,
        model: ollamaModel
      });
      
      if (result) {
        setMessage({type: 'success', text: 'Configuração do Ollama é válida!'});
      } else {
        setMessage({type: 'error', text: 'Configuração do Ollama é inválida!'});
      }
    } catch (err: any) {
      setMessage({type: 'error', text: `Erro ao testar configuração: ${err.message}`});
    } finally {
      setIsLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="config-modal-overlay">
      <div className="config-modal">
        <div className="config-modal-header">
          <h2>Configurações</h2>
          <button className="config-modal-close" onClick={onClose}>×</button>
        </div>
        
        <div className="config-modal-body">
          {/* OpenAI Configuration */}
          <div className="config-section">
            <h3>OpenAI API</h3>
            <div className="config-field">
              <label htmlFor="openai-key">Chave da API:</label>
              <input
                id="openai-key"
                type="password"
                value={openAIKey}
                onChange={(e) => setOpenAIKey(e.target.value)}
                placeholder="sk-..."
              />
              <button 
                className="btn test-btn" 
                onClick={handleTestOpenAI}
                disabled={isLoading || !openAIKey}
              >
                Testar
              </button>
            </div>
          </div>
          
          {/* Ollama Configuration */}
          <div className="config-section">
            <h3>Ollama</h3>
            <div className="config-field">
              <label htmlFor="ollama-url">URL Base:</label>
              <input
                id="ollama-url"
                type="text"
                value={ollamaURL}
                onChange={(e) => setOllamaURL(e.target.value)}
                placeholder="http://localhost:11434"
              />
            </div>
            <div className="config-field">
              <label htmlFor="ollama-model">Modelo:</label>
              <input
                id="ollama-model"
                type="text"
                value={ollamaModel}
                onChange={(e) => setOllamaModel(e.target.value)}
                placeholder="llama2"
              />
              <button 
                className="btn test-btn" 
                onClick={handleTestOllama}
                disabled={isLoading || (!ollamaURL && !ollamaModel)}
              >
                Testar
              </button>
            </div>
          </div>
          
          {/* Message Display */}
          {message && (
            <div className={`message message-${message.type}`}>
              {message.text}
            </div>
          )}
        </div>
        
        <div className="config-modal-footer">
          <button className="btn cancel-btn" onClick={onClose} disabled={isLoading}>
            Cancelar
          </button>
          <button className="btn save-btn" onClick={handleSave} disabled={isLoading}>
            {isLoading ? 'Salvando...' : 'Salvar'}
          </button>
        </div>
      </div>
    </div>
  );
}

export default ConfigModal;