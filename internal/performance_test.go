// Package internal tests de performance e stress testing
package internal

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"

	"BrxAgente-desafio4/internal/cache"
	"BrxAgente-desafio4/internal/modelo"
	"BrxAgente-desafio4/internal/monitoring"
	"BrxAgente-desafio4/internal/parallel"
)

// TestStressLLMCache testa o cache LLM sob stress
func TestStressLLMCache(t *testing.T) {
	llmCache := cache.NewLLMCache(1000, time.Hour)
	numWorkers := 50
	numOperations := 10000

	var wg sync.WaitGroup
	start := time.Now()

	// Workers fazendo operações concorrentes no cache
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			
			for j := 0; j < numOperations/numWorkers; j++ {
				prompt := fmt.Sprintf("worker_%d_prompt_%d", workerID, j)
				response := fmt.Sprintf("response_for_worker_%d_op_%d", workerID, j)
				
				metadata := cache.CacheMetadata{
					TokensUsed: 100,
					Duration:   time.Millisecond * 50,
					Model:      "test-model",
					Quality:    0.95,
				}
				
				// Set operation
				llmCache.Set(prompt, response, metadata)
				
				// Get operation
				if got, found := llmCache.Get(prompt); !found || got != response {
					t.Errorf("Cache miss ou valor incorreto para worker %d op %d", workerID, j)
				}
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	stats := llmCache.GetStats()
	
	t.Logf("Stress Test LLM Cache:")
	t.Logf("- Tempo total: %v", elapsed)
	t.Logf("- Operações: %d", numOperations)
	t.Logf("- Workers: %d", numWorkers)
	t.Logf("- Cache size: %d", stats.Size)
	t.Logf("- Hit ratio: %.2f%%", stats.HitRatio*100)
	t.Logf("- Operações/segundo: %.2f", float64(numOperations)/elapsed.Seconds())

	if stats.HitRatio < 0.8 {
		t.Errorf("Hit ratio muito baixo: %.2f", stats.HitRatio)
	}
}

// TestStressDataCache testa o cache de dados sob stress
func TestStressDataCache(t *testing.T) {
	dataCache := cache.NewDataCache(50*1024*1024, time.Hour) // 50MB
	numColaboradores := 10000

	// Cria dataset grande de colaboradores
	colaboradores := make(map[string]*modelo.Colaborador, numColaboradores)
	for i := 0; i < numColaboradores; i++ {
		colaboradores[fmt.Sprintf("MAT%05d", i)] = &modelo.Colaborador{
			Matricula: fmt.Sprintf("MAT%05d", i),
			Nome:      fmt.Sprintf("Colaborador %d", i),
			Sindicato: fmt.Sprintf("SIND%02d", i%10),
		}
	}

	start := time.Now()
	
	// Armazena no cache
	key := "colaboradores_stress_test"
	dataCache.SetColaboradores(key, colaboradores)
	
	// Recupera do cache múltiplas vezes
	for i := 0; i < 100; i++ {
		if cached, found := dataCache.GetColaboradores(key); !found {
			t.Errorf("Falha recuperando colaboradores do cache")
		} else if len(cached) != numColaboradores {
			t.Errorf("Número incorreto de colaboradores: esperado %d, obtido %d", numColaboradores, len(cached))
		}
	}

	elapsed := time.Since(start)
	stats := dataCache.GetStats()

	t.Logf("Stress Test Data Cache:")
	t.Logf("- Tempo total: %v", elapsed)
	t.Logf("- Colaboradores: %d", numColaboradores)
	t.Logf("- Cache entries: %d", stats.EntryCount)
	t.Logf("- Uso de memória: %.2f MB", float64(stats.CurrentSize)/1024/1024)
	t.Logf("- Uso percentual: %.2f%%", stats.UsagePercent)
}

// BenchmarkWorkerPool benchmark do worker pool
func BenchmarkWorkerPool(b *testing.B) {
	sizes := []int{1, 10, 50, 100}
	
	for _, workers := range sizes {
		b.Run(fmt.Sprintf("workers_%d", workers), func(b *testing.B) {
			pool := parallel.NewWorkerPool(workers, 1000)
			defer pool.Close()

			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				task := &TestTask{
					ID:       fmt.Sprintf("task_%d", i),
					Duration: time.Millisecond,
				}
				
				pool.SubmitTask(task)
			}
			
			pool.WaitForCompletion(30 * time.Second)
		})
	}
}

// BenchmarkBatchProcessor benchmark do processador de lotes
func BenchmarkBatchProcessor(b *testing.B) {
	config := parallel.BatchConfig{
		WorkerPoolSize: 10,
		BatchSize:      100,
		MaxConcurrent:  5,
		TimeoutSeconds: 30,
	}
	
	processor := parallel.NewBatchProcessor(config)
	defer processor.Close()

	// Cria colaboradores de teste
	colaboradores := make(map[string]*modelo.Colaborador)
	for i := 0; i < 1000; i++ {
		matricula := fmt.Sprintf("MAT%05d", i)
		colaboradores[matricula] = &modelo.Colaborador{
			Matricula: matricula,
			Nome:      fmt.Sprintf("Colaborador %d", i),
			Sindicato: "SIND01",
		}
	}

	// Função de processamento simples
	processFunc := func(col *modelo.Colaborador) (*modelo.Colaborador, error) {
		// Simula processamento
		time.Sleep(time.Microsecond * 100)
		return col, nil
	}

	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, erros := processor.ProcessColaboradoresBatch(colaboradores, processFunc)
		if len(erros) > 0 {
			b.Errorf("Erros no processamento: %v", erros)
		}
	}
}

// TestLargeVolumeProcessing testa processamento de grande volume
func TestLargeVolumeProcessing(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste de grande volume em modo short")
	}

	numColaboradores := 10000
	t.Logf("Iniciando teste de volume com %d colaboradores", numColaboradores)

	// Cria dataset grande
	colaboradores := make(map[string]*modelo.Colaborador)
	for i := 0; i < numColaboradores; i++ {
		matricula := fmt.Sprintf("MAT%06d", i)
		colaboradores[matricula] = &modelo.Colaborador{
			Matricula:    matricula,
			Nome:         fmt.Sprintf("Colaborador Teste %d", i),
			Sindicato:    fmt.Sprintf("SIND%02d", i%20),
			DataAdmissao: time.Now().AddDate(-rand.Intn(5), -rand.Intn(12), -rand.Intn(30)),
		}
	}

	// Configuração otimizada para grande volume
	config := parallel.BatchConfig{
		WorkerPoolSize: runtime.NumCPU() * 2,
		BatchSize:      200,
		MaxConcurrent:  runtime.NumCPU(),
		TimeoutSeconds: 60,
	}

	processor := parallel.NewBatchProcessor(config)
	defer processor.Close()

	// Métricas de monitoramento
	metrics := monitoring.NewMetricsCollector(1)
	defer metrics.Stop()

	start := time.Now()

	// Processa lote grande
	processFunc := func(col *modelo.Colaborador) (*modelo.Colaborador, error) {
		// Simula processamento complexo
		time.Sleep(time.Microsecond * time.Duration(100+rand.Intn(200)))
		
		// Simula erro ocasional
		if rand.Float64() < 0.01 { // 1% de erro
			return nil, fmt.Errorf("erro simulado para %s", col.Matricula)
		}
		
		return col, nil
	}

	resultados, erros := processor.ProcessColaboradoresBatch(colaboradores, processFunc)
	
	elapsed := time.Since(start)
	throughput := float64(numColaboradores) / elapsed.Seconds()

	t.Logf("Teste de Grande Volume - Resultados:")
	t.Logf("- Colaboradores processados: %d", len(resultados))
	t.Logf("- Erros: %d", len(erros))
	t.Logf("- Tempo total: %v", elapsed)
	t.Logf("- Throughput: %.2f colaboradores/segundo", throughput)

	stats := processor.GetStats()
	t.Logf("- Workers utilizados: %d", stats.WorkerPoolStats.ActiveWorkers)
	t.Logf("- Tarefas processadas: %d", stats.WorkerPoolStats.TasksProcessed)
	t.Logf("- Tempo médio por tarefa: %v", stats.WorkerPoolStats.AverageDuration)

	// Verifica targets de performance
	if elapsed > 5*time.Minute {
		t.Errorf("Processamento muito lento: %v > 5min", elapsed)
	}

	if throughput < 50 { // Mínimo 50 colaboradores/segundo
		t.Errorf("Throughput muito baixo: %.2f < 50/sec", throughput)
	}

	errorRate := float64(len(erros)) / float64(numColaboradores) * 100
	if errorRate > 2 { // Máximo 2% de erro
		t.Errorf("Taxa de erro muito alta: %.2f%% > 2%%", errorRate)
	}
}

// TestMemoryUsageUnderLoad testa uso de memória sob carga
func TestMemoryUsageUnderLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste de memória em modo short")
	}

	metrics := monitoring.NewMetricsCollector(1)
	defer metrics.Stop()

	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// Simula carga pesada
	smartCache := cache.NewSmartCache(cache.SmartCacheConfig{
		LLMMaxSize:       1000,
		LLMTTLHours:     1,
		DataMaxSizeBytes: 100 * 1024 * 1024, // 100MB
		DataTTLHours:    1,
	})

	// Cria muitos dados em cache
	for i := 0; i < 5000; i++ {
		prompt := fmt.Sprintf("prompt_de_teste_numero_%d_com_texto_longo", i)
		response := fmt.Sprintf("resposta_detalhada_para_prompt_%d_com_muito_conteudo", i)
		
		metadata := cache.CacheMetadata{
			TokensUsed: 200 + rand.Intn(300),
			Duration:   time.Millisecond * time.Duration(50+rand.Intn(200)),
			Model:      "test-model-large",
			Quality:    0.8 + rand.Float64()*0.2,
		}

		smartCache.SetLLMResponse(prompt, response, metadata)

		// Adiciona dados no cache de dados também
		if i%100 == 0 {
			colaboradores := make(map[string]*modelo.Colaborador)
			for j := 0; j < 100; j++ {
				matricula := fmt.Sprintf("MAT_%d_%d", i, j)
				colaboradores[matricula] = &modelo.Colaborador{
					Matricula: matricula,
					Nome:      fmt.Sprintf("Nome %d %d", i, j),
					Sindicato: "SIND_TEST",
				}
			}
			
			key := fmt.Sprintf("colaboradores_batch_%d", i/100)
			smartCache.SetProcessedData(key, colaboradores, 20*1024) // 20KB estimado
		}
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)

	memoryUsedMB := float64(m2.Alloc-m1.Alloc) / 1024 / 1024
	stats := smartCache.GetCombinedStats()

	t.Logf("Teste de Uso de Memória:")
	t.Logf("- Memória usada: %.2f MB", memoryUsedMB)
	t.Logf("- LLM Cache size: %d", stats.LLMStats.Size)
	t.Logf("- LLM Cache hit ratio: %.2f%%", stats.LLMStats.HitRatio*100)
	t.Logf("- Data Cache entries: %d", stats.DataStats.EntryCount)
	t.Logf("- Data Cache size: %.2f MB", float64(stats.DataStats.CurrentSize)/1024/1024)

	// Verifica limites de memória
	if memoryUsedMB > 1000 { // 1GB limite
		t.Errorf("Uso de memória excessivo: %.2f MB > 1000 MB", memoryUsedMB)
	}

	// Testa otimização de cache
	optimizationResult := smartCache.OptimizeCaches()
	if !optimizationResult.Success {
		t.Errorf("Falha na otimização de cache")
	}

	t.Logf("Otimização de Cache:")
	t.Logf("- LLM entries removidas: %d", optimizationResult.LLMEntriesRemoved)
	t.Logf("- Data entries removidas: %d", optimizationResult.DataEntriesRemoved)
}

// TestConcurrentCacheOperations testa operações concorrentes no cache
func TestConcurrentCacheOperations(t *testing.T) {
	smartCache := cache.NewSmartCache(cache.SmartCacheConfig{
		LLMMaxSize:       500,
		LLMTTLHours:     1,
		DataMaxSizeBytes: 50 * 1024 * 1024,
		DataTTLHours:    1,
	})

	numGoroutines := 100
	operationsPerGoroutine := 1000

	var wg sync.WaitGroup
	start := time.Now()

	// Múltiplas goroutines fazendo operações concorrentes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < operationsPerGoroutine; j++ {
				// Operações LLM Cache
				prompt := fmt.Sprintf("g%d_prompt_%d", goroutineID, j)
				response := fmt.Sprintf("response_g%d_op%d", goroutineID, j)
				
				metadata := cache.CacheMetadata{
					TokensUsed: 100,
					Duration:   time.Millisecond * 50,
					Model:      "test",
					Quality:    0.9,
				}

				smartCache.SetLLMResponse(prompt, response, metadata)

				if got, found := smartCache.GetLLMResponse(prompt); !found || got != response {
					t.Errorf("Cache inconsistente para goroutine %d op %d", goroutineID, j)
					return
				}

				// Operações Data Cache a cada 50 operações
				if j%50 == 0 {
					key := fmt.Sprintf("data_g%d_op%d", goroutineID, j)
					data := map[string]interface{}{
						"goroutine": goroutineID,
						"operation": j,
						"timestamp": time.Now(),
					}

					smartCache.SetProcessedData(key, data, 1024)

					if _, found := smartCache.GetProcessedData(key); !found {
						t.Errorf("Data cache miss para goroutine %d op %d", goroutineID, j)
						return
					}
				}
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)
	totalOps := numGoroutines * operationsPerGoroutine

	stats := smartCache.GetCombinedStats()

	t.Logf("Teste de Concorrência:")
	t.Logf("- Operações totais: %d", totalOps)
	t.Logf("- Tempo total: %v", elapsed)
	t.Logf("- Operações/segundo: %.2f", float64(totalOps)/elapsed.Seconds())
	t.Logf("- LLM Hit ratio: %.2f%%", stats.LLMStats.HitRatio*100)
	t.Logf("- Data cache entries: %d", stats.DataStats.EntryCount)

	if stats.LLMStats.HitRatio < 0.9 {
		t.Errorf("Hit ratio baixo em teste de concorrência: %.2f", stats.LLMStats.HitRatio)
	}
}

// TestTask implementa Task interface para testes
type TestTask struct {
	ID       string
	Duration time.Duration
	ShouldError bool
}

func (tt *TestTask) Execute(ctx context.Context) (parallel.Result, error) {
	if tt.ShouldError {
		return nil, fmt.Errorf("erro simulado para task %s", tt.ID)
	}

	time.Sleep(tt.Duration)

	return &TestResult{
		ID:       tt.ID,
		Data:     fmt.Sprintf("resultado_para_%s", tt.ID),
		Duration: tt.Duration,
	}, nil
}

func (tt *TestTask) GetID() string {
	return tt.ID
}

func (tt *TestTask) GetPriority() int {
	return 0
}

// TestResult implementa Result interface para testes
type TestResult struct {
	ID       string
	Data     interface{}
	Duration time.Duration
}

func (tr *TestResult) GetID() string {
	return tr.ID
}

func (tr *TestResult) GetData() interface{} {
	return tr.Data
}

func (tr *TestResult) GetDuration() time.Duration {
	return tr.Duration
}