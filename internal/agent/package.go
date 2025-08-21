// Package agent implementa o agente inteligente baseado em LangChainGo
// para automação de processos de Vale Refeição (VR).
//
// Este pacote fornece funcionalidades para:
//   - Execução de workflows automatizados
//   - Processamento inteligente de dados de VR
//   - Interação em linguagem natural
//   - Integração com ferramentas específicas do domínio
//
// Exemplo de uso básico:
//
//	cfg := config.LoadAgentConfig()
//	agent := agent.NewVRAgent(cfg, chatService, excelService, calculoService)
//	response, err := agent.Ask("Quantos colaboradores ativos temos?")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(response)
//
// O agente utiliza LangChainGo para orquestrar chamadas a LLMs (OpenAI ou Ollama)
// e integra com os sistemas existentes através de ferramentas especializadas.
package agent