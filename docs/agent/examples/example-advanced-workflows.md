# 🔧 Exemplo: Workflows Avançados - Personalização e Automação

Este exemplo demonstra como criar e personalizar workflows avançados no BrxAgente para cenários específicos da sua organização.

## 📋 Cenário

**Situação:** Sua empresa tem necessidades específicas de processamento que não são cobertas pelos workflows padrão:
- Processamento apenas de determinados sindicatos
- Validação rigorosa com critérios customizados  
- Geração de relatórios personalizados
- Integração com sistemas externos

**Objetivo:** Criar workflows personalizados e automatizar rotinas mensais.

## 🎯 Workflow 1: Processamento por Sindicato Específico

### **Configuração do Workflow**

1. **Acessar Configurações de Workflow**
   - Menu: **"⚙️ Configurações"**
   - Aba: **"🔄 Workflows"**
   - Clique em **"➕ Novo Workflow"**

2. **Definir Parâmetros Base**
   ```yaml
   nome: processamento-sindpd-apenas
   descricao: Processa apenas colaboradores do SINDPD
   tipo: processamento-filtrado
   
   parametros:
     sindicatos_incluir: ["SINDPD"]
     sindicatos_excluir: []
     incluir_inativos: false
     validacao_rigorosa: true
   
   etapas:
     - validacao_diretorio
     - leitura_planilhas
     - filtragem_sindicatos
     - consolidacao_dados
     - aplicacao_regras
     - geracao_planilha
     - relatorio_final
   ```

3. **Configurar Filtros Avançados**
   ```yaml
   filtros:
     por_sindicato:
       ativo: true
       lista: ["SINDPD"]
       modo: inclusivo
     
     por_status:
       incluir_ativos: true
       incluir_inativos: false
       incluir_afastados: true
     
     por_data_admissao:
       ativo: false
       data_inicio: "2020-01-01"
       data_fim: "2025-12-31"
     
     por_cargo:
       excluir_diretores: true
       excluir_estagiarios: true
       excluir_aprendizes: true
   ```

### **Execução do Workflow**

1. **Iniciar Processamento**
   - Dashboard: **"🔄 Workflows"**
   - Dropdown: **"Processamento SINDPD Apenas"**
   - Configurar pasta das planilhas
   - Clique em **"▶️ Executar"**

2. **Acompanhar Progresso**
   ```
   🔄 Processamento SINDPD - Dezembro 2025
   ██████████████████████░░ 90%
   
   📋 Etapas Concluídas:
   1. ✅ Validação do diretório      (12s)
   2. ✅ Leitura das planilhas       (1m 45s)
   3. ✅ Filtragem por sindicato     (23s)
   4. ✅ Consolidação dos dados      (2m 12s)
   5. ✅ Aplicação das regras        (3m 45s)
   6. 🔄 Geração da planilha         (1m 32s...)
   
   📊 Dados Filtrados:
   • Total de colaboradores: 2.847
   • Após filtro SINDPD: 1.247
   • Colaboradores ativos: 1.189
   • Colaboradores afastados: 58
   ```

3. **Resultado Final**
   ```
   ✅ PROCESSAMENTO SINDPD CONCLUÍDO
   
   👥 Colaboradores processados: 1.247 (apenas SINDPD)
   💰 VR total calculado: R$ 573.620,00
   📊 Média por colaborador: R$ 460,11
   ⏱️  Tempo total: 9 minutos e 57 segundos
   
   📁 Arquivos gerados:
   • output/vr-sindpd-dezembro-2025.xlsx
   • output/relatorio-sindpd-dezembro.xlsx  
   • output/log-processamento-sindpd.txt
   ```

## 🔍 Workflow 2: Validação Rigorosa com Alertas

### **Configuração de Validação Avançada**

1. **Criar Workflow de Validação**
   ```yaml
   nome: validacao-rigorosa-alertas
   descricao: Validação completa com alertas automáticos
   tipo: validacao-avancada
   
   parametros:
     nivel_validacao: rigoroso
     gerar_alertas: true
     parar_em_erro_critico: true
     salvar_log_detalhado: true
   
   validacoes:
     estrutura_arquivos:
       verificar_colunas_obrigatorias: true
       verificar_formatos_data: true
       verificar_tipos_dados: true
       
     consistencia_dados:
       verificar_matriculas_duplicadas: true
       verificar_relacionamentos: true
       verificar_datas_logicas: true
       
     regras_negocio:
       verificar_sindicatos_validos: true
       verificar_tipos_afastamento: true
       verificar_periodos_sobrepostos: true
   ```

2. **Configurar Critérios de Alerta**
   ```yaml
   alertas:
     email:
       ativo: true
       destinatarios: ["rh@empresa.com", "ti@empresa.com"]
       nivel_minimo: warning
       
     dashboard:
       ativo: true
       notificacoes_popup: true
       sons_alerta: false
       
     arquivo_log:
       ativo: true
       nivel_detalhe: completo
       rotacao_logs: mensal
   
   criterios_alerta:
     matriculas_duplicadas:
       nivel: error
       acao: parar_processamento
       
     relacionamentos_quebrados:
       nivel: warning
       acao: continuar_com_log
       
     datas_inconsistentes:
       nivel: error
       acao: solicitar_correcao
   ```

### **Execução com Alertas**

1. **Iniciar Validação**
   - Selecionar workflow **"Validação Rigorosa com Alertas"**
   - Configurar pasta de planilhas
   - Marcar **"📧 Enviar alertas por email"**
   - Executar

2. **Exemplo de Alertas Gerados**
   ```
   🚨 ALERTAS DE VALIDAÇÃO DETECTADOS
   
   ❌ ERROS CRÍTICOS (2):
   1. colaboradores.xlsx - Linha 1847
      • Matrícula duplicada: MAT003456
      • Ação: Processamento interrompido
      
   2. afastamentos.xlsx - Linha 234
      • Data inválida: 30/02/2025
      • Ação: Registro ignorado
   
   ⚠️  AVISOS (7):
   1. colaboradores.xlsx - Linha 892
      • Sindicato não reconhecido: "SINDX"
      • Sugestão: Verificar se é "SINDPD"
      
   2. afastamentos.xlsx - Linha 156
      • Afastamento sem colaborador: MAT999999
      • Ação: Registro será ignorado
   
   📧 Email de alerta enviado para:
   • rh@empresa.com
   • ti@empresa.com
   
   📋 Ações Recomendadas:
   [ ] Corrigir matrícula duplicada na linha 1847
   [ ] Corrigir data inválida 30/02/2025
   [ ] Revisar sindicato "SINDX" - linha 892
   [ ] Verificar matrícula MAT999999 em afastamentos
   ```

## 📊 Workflow 3: Geração de Relatórios Executivos

### **Configuração do Workflow**

```yaml
nome: relatorio-executivo-mensal
descricao: Gera relatório completo para diretoria
tipo: relatorio-executivo

parametros:
  incluir_graficos: true
  incluir_comparativo_historico: true
  incluir_predicoes: true
  incluir_recomendacoes: true
  formato_saida: ["xlsx", "pdf"]

secoes_relatorio:
  resumo_executivo:
    ativo: true
    incluir_metricas_chave: true
    incluir_alertas_importantes: true
    
  analise_por_sindicato:
    ativo: true
    incluir_breakdown_custos: true
    incluir_comparativo_mes_anterior: true
    
  deteccao_anomalias:
    ativo: true
    incluir_analise_causa: true
    incluir_acao_recomendada: true
    
  predicoes_tendencias:
    ativo: true
    horizonte_predicao: 6 # meses
    incluir_cenarios: ["otimista", "realista", "pessimista"]
```

### **Exemplo de Relatório Gerado**

```
📊 RELATÓRIO EXECUTIVO VR - DEZEMBRO 2025

═══════════════════════════════════════
🎯 RESUMO EXECUTIVO
═══════════════════════════════════════

💰 INVESTIMENTO TOTAL: R$ 1.309.620,00
👥 COLABORADORES BENEFICIADOS: 2.847
📈 VARIAÇÃO MÊS ANTERIOR: +2.3% (+R$ 29.840,00)
⚠️  ALERTAS IMPORTANTES: 3 anomalias detectadas

═══════════════════════════════════════
📊 BREAKDOWN POR SINDICATO
═══════════════════════════════════════

🏢 SINDPD (43.8% do total)
• Colaboradores: 1.247 (+12 vs Nov/25)
• Investimento: R$ 573.620,00
• Média por colaborador: R$ 460,11
• Status: ✅ Dentro do esperado

🏭 SINDAC (31.3% do total)  
• Colaboradores: 892 (-3 vs Nov/25)
• Investimento: R$ 410.320,00
• Média por colaborador: R$ 459,91
• Status: ✅ Dentro do esperado

⚙️ SINDMET (24.9% do total)
• Colaboradores: 708 (+8 vs Nov/25)
• Investimento: R$ 325.680,00
• Média por colaborador: R$ 459,84
• Status: ✅ Dentro do esperado

═══════════════════════════════════════
🚨 ANOMALIAS DETECTADAS
═══════════════════════════════════════

1. 🔍 VALOR ALTO DETECTADO
   • Colaborador: MAT003789
   • VR Calculado: R$ 920,00 (200% acima da média)
   • Causa Provável: Erro de cálculo ou caso especial
   • Ação: ✅ Corrigido automaticamente para R$ 460,00

2. ⚠️  PADRÃO INCOMUM - SINDAC
   • Redução de 15 colaboradores em relação à média
   • Pode indicar: desligamentos não reportados
   • Ação: Verificação manual recomendada

3. 📅 DATAS INCONSISTENTES
   • 2 registros com datas de afastamento sobrepostas
   • Status: ✅ Corrigido automaticamente

═══════════════════════════════════════
🔮 PREDIÇÕES E TENDÊNCIAS
═══════════════════════════════════════

📈 TENDÊNCIA PRÓXIMOS 6 MESES:
• Janeiro 2026: R$ 1.325.000,00 (±3%)
• Fevereiro 2026: R$ 1.341.000,00 (±3%)
• Março 2026: R$ 1.358.000,00 (±4%)

💡 RECOMENDAÇÕES ESTRATÉGICAS:
1. Considerar revisão de valores SINDPD (+2.1% sugerido)
2. Implementar validação prévia para anomalias recorrentes
3. Avaliar impacto de crescimento projetado no orçamento Q1/2026

═══════════════════════════════════════
📋 PRÓXIMAS AÇÕES
═══════════════════════════════════════

🎯 IMEDIATAS (até 5 dias):
[ ] Revisar anomalia SINDAC (15 colaboradores)
[ ] Validar correções aplicadas automaticamente  
[ ] Confirmar orçamento janeiro com Financeiro

📅 MÊS QUE VEM:
[ ] Executar processamento janeiro até dia 25
[ ] Implementar melhorias sugeridas pela IA
[ ] Revisar acordos sindicais para 2026
```

## 🔄 Workflow 4: Automação Completa com Agendamento

### **Configuração de Automação**

1. **Criar Workflow Automatizado**
   ```yaml
   nome: processamento-automatico-mensal
   descricao: Execução automática todo mês
   tipo: automatizado
   
   agendamento:
     frequencia: mensal
     dia_execucao: 25  # Todo dia 25
     horario: "08:00"
     fuso_horario: "America/Sao_Paulo"
     
   pre_requisitos:
     verificar_planilhas_atualizadas: true
     verificar_conexao_api: true
     verificar_espaco_disco: true
     
   pos_processamento:
     enviar_relatorio_email: true
     fazer_backup_resultados: true
     limpar_logs_antigos: true
     atualizar_dashboard: true
   ```

2. **Configurar Monitoramento**
   ```yaml
   monitoramento:
     alertas_tempo_execucao:
       ativo: true
       tempo_maximo: 30 # minutos
       acao_timeout: enviar_alerta
       
     alertas_qualidade:
       ativo: true
       anomalias_maximas: 5
       confiabilidade_minima: 95%
       
     notificacoes_sucesso:
       email: ["rh@empresa.com"]
       slack: "#rh-notifications"
       dashboard: true
   ```

### **Dashboard de Automação**

```
🤖 AUTOMAÇÃO VR - PAINEL DE CONTROLE

═══════════════════════════════════════
⏰ PRÓXIMA EXECUÇÃO: 25/01/2026 às 08:00
═══════════════════════════════════════

📅 HISTÓRICO (ÚLTIMOS 6 MESES):
• Dezembro/25: ✅ Sucesso (9m 32s) - 2.847 colaboradores
• Novembro/25: ✅ Sucesso (8m 45s) - 2.834 colaboradores  
• Outubro/25: ✅ Sucesso (9m 12s) - 2.821 colaboradores
• Setembro/25: ⚠️  Sucesso com avisos (11m 23s) - 2.809 colaboradores
• Agosto/25: ✅ Sucesso (8m 56s) - 2.798 colaboradores
• Julho/25: ❌ Falhou - Planilhas não encontradas

📊 ESTATÍSTICAS:
• Taxa de sucesso: 83.3% (5/6)
• Tempo médio: 9m 35s
• Colaboradores/mês (média): 2.822

🔧 CONFIGURAÇÕES ATIVAS:
• [✓] Validação automática
• [✓] Correção de anomalias simples  
• [✓] Backup automático
• [✓] Relatório executivo
• [✓] Notificações email
• [ ] Integração ERP (em desenvolvimento)

⚙️ AÇÕES:
[📝 Editar Agendamento] [🧪 Teste Manual] [📊 Ver Logs] [⚠️ Pausar Automação]
```

## 💡 Dicas de Otimização

### **Performance**
- Use filtros para processar apenas dados necessários
- Configure cache para melhorar velocidade
- Execute workflows pesados fora do horário comercial

### **Qualidade**
- Sempre use validação rigorosa em produção
- Configure alertas para anomalias críticas
- Mantenha logs detalhados para auditoria

### **Manutenção**
- Revise workflows mensalmente
- Atualize critérios conforme mudanças organizacionais
- Teste mudanças em ambiente separado

---

**🚀 Com esses workflows avançados, você pode automatizar completamente o processamento de VR e ter insights detalhados sobre os dados da sua organização!**