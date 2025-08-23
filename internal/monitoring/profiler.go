package monitoring

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

// AutoProfiler gerencia profiling automático do sistema
type AutoProfiler struct {
	enabled      bool
	outputDir    string
	cpuFile      *os.File
	memInterval  time.Duration
	stopCh       chan bool
	mutex        sync.Mutex
	isRunning    bool
	profileCount int
}

// ProfilerConfig configuração do profiler automático
type ProfilerConfig struct {
	Enabled             bool
	OutputDir           string
	MemProfileInterval  time.Duration // Intervalo para profile de memória
	AutoCPUThreshold    float64       // Threshold de CPU para iniciar profiling
	AutoMemThreshold    float64       // Threshold de memória (MB) para profile
}

// NewAutoProfiler cria um novo profiler automático
func NewAutoProfiler(config ProfilerConfig) *AutoProfiler {
	if config.OutputDir == "" {
		config.OutputDir = "./profiles"
	}

	if config.MemProfileInterval == 0 {
		config.MemProfileInterval = 10 * time.Minute
	}

	// Cria diretório se não existir
	os.MkdirAll(config.OutputDir, 0755)

	return &AutoProfiler{
		enabled:     config.Enabled,
		outputDir:   config.OutputDir,
		memInterval: config.MemProfileInterval,
		stopCh:      make(chan bool),
	}
}

// StartContinuousProfiling inicia profiling contínuo
func (ap *AutoProfiler) StartContinuousProfiling() error {
	if !ap.enabled {
		return fmt.Errorf("profiler não está habilitado")
	}

	ap.mutex.Lock()
	defer ap.mutex.Unlock()

	if ap.isRunning {
		return fmt.Errorf("profiler já está executando")
	}

	ap.isRunning = true

	// Inicia goroutine para profiles de memória periódicos
	go ap.memoryProfileLoop()

	return nil
}

// StartCPUProfile inicia profiling de CPU
func (ap *AutoProfiler) StartCPUProfile() error {
	if !ap.enabled {
		return nil
	}

	ap.mutex.Lock()
	defer ap.mutex.Unlock()

	if ap.cpuFile != nil {
		return fmt.Errorf("CPU profiling já está ativo")
	}

	filename := filepath.Join(ap.outputDir, fmt.Sprintf("cpu_%d.prof", time.Now().Unix()))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("erro criando arquivo de CPU profile: %w", err)
	}

	ap.cpuFile = file
	return pprof.StartCPUProfile(file)
}

// StopCPUProfile para profiling de CPU
func (ap *AutoProfiler) StopCPUProfile() error {
	if !ap.enabled {
		return nil
	}

	ap.mutex.Lock()
	defer ap.mutex.Unlock()

	if ap.cpuFile == nil {
		return nil // Não estava ativo
	}

	pprof.StopCPUProfile()
	ap.cpuFile.Close()
	ap.cpuFile = nil

	return nil
}

// CreateMemoryProfile cria um profile de memória
func (ap *AutoProfiler) CreateMemoryProfile() error {
	if !ap.enabled {
		return nil
	}

	ap.mutex.Lock()
	ap.profileCount++
	count := ap.profileCount
	ap.mutex.Unlock()

	filename := filepath.Join(ap.outputDir, fmt.Sprintf("mem_%d_%d.prof", time.Now().Unix(), count))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("erro criando arquivo de memory profile: %w", err)
	}
	defer file.Close()

	runtime.GC() // Força garbage collection antes do profile

	if err := pprof.WriteHeapProfile(file); err != nil {
		return fmt.Errorf("erro escrevendo memory profile: %w", err)
	}

	return nil
}

// CreateGoroutineProfile cria um profile de goroutines
func (ap *AutoProfiler) CreateGoroutineProfile() error {
	if !ap.enabled {
		return nil
	}

	filename := filepath.Join(ap.outputDir, fmt.Sprintf("goroutine_%d.prof", time.Now().Unix()))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("erro criando arquivo de goroutine profile: %w", err)
	}
	defer file.Close()

	if err := pprof.Lookup("goroutine").WriteTo(file, 1); err != nil {
		return fmt.Errorf("erro escrevendo goroutine profile: %w", err)
	}

	return nil
}

// CreateBlockProfile cria um profile de blocking
func (ap *AutoProfiler) CreateBlockProfile() error {
	if !ap.enabled {
		return nil
	}

	filename := filepath.Join(ap.outputDir, fmt.Sprintf("block_%d.prof", time.Now().Unix()))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("erro criando arquivo de block profile: %w", err)
	}
	defer file.Close()

	if err := pprof.Lookup("block").WriteTo(file, 1); err != nil {
		return fmt.Errorf("erro escrevendo block profile: %w", err)
	}

	return nil
}

// CreateFullSnapshot cria um snapshot completo
func (ap *AutoProfiler) CreateFullSnapshot() error {
	if !ap.enabled {
		return nil
	}

	timestamp := time.Now().Unix()
	
	// Cria diretório para o snapshot
	snapshotDir := filepath.Join(ap.outputDir, fmt.Sprintf("snapshot_%d", timestamp))
	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		return fmt.Errorf("erro criando diretório de snapshot: %w", err)
	}

	// Profile de memória
	if err := ap.createSnapshotProfile(snapshotDir, "heap", "heap"); err != nil {
		return err
	}

	// Profile de goroutines
	if err := ap.createSnapshotProfile(snapshotDir, "goroutine", "goroutine"); err != nil {
		return err
	}

	// Profile de block
	if err := ap.createSnapshotProfile(snapshotDir, "block", "block"); err != nil {
		return err
	}

	// Profile de mutex
	if err := ap.createSnapshotProfile(snapshotDir, "mutex", "mutex"); err != nil {
		return err
	}

	return nil
}

// createSnapshotProfile cria um profile específico no snapshot
func (ap *AutoProfiler) createSnapshotProfile(dir, profileType, filename string) error {
	filepath := filepath.Join(dir, filename+".prof")
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("erro criando arquivo %s: %w", filename, err)
	}
	defer file.Close()

	profile := pprof.Lookup(profileType)
	if profile == nil {
		return fmt.Errorf("profile %s não encontrado", profileType)
	}

	if err := profile.WriteTo(file, 1); err != nil {
		return fmt.Errorf("erro escrevendo profile %s: %w", profileType, err)
	}

	return nil
}

// Stop para o profiler automático
func (ap *AutoProfiler) Stop() error {
	if !ap.enabled {
		return nil
	}

	ap.mutex.Lock()
	defer ap.mutex.Unlock()

	if !ap.isRunning {
		return nil
	}

	// Para CPU profiling se estiver ativo
	if ap.cpuFile != nil {
		pprof.StopCPUProfile()
		ap.cpuFile.Close()
		ap.cpuFile = nil
	}

	// Sinaliza para parar loops
	close(ap.stopCh)
	ap.isRunning = false

	return nil
}

// memoryProfileLoop executa loop para profiles de memória periódicos
func (ap *AutoProfiler) memoryProfileLoop() {
	ticker := time.NewTicker(ap.memInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := ap.CreateMemoryProfile(); err != nil {
				// Log do erro (em produção usaria logger proper)
				fmt.Printf("Erro criando memory profile: %v\n", err)
			}

		case <-ap.stopCh:
			return
		}
	}
}

// GetProfilerStats retorna estatísticas do profiler
func (ap *AutoProfiler) GetProfilerStats() ProfilerStats {
	ap.mutex.Lock()
	defer ap.mutex.Unlock()

	stats := ProfilerStats{
		Enabled:          ap.enabled,
		IsRunning:        ap.isRunning,
		OutputDirectory:  ap.outputDir,
		ProfilesCreated:  ap.profileCount,
		CPUProfileActive: ap.cpuFile != nil,
		MemoryInterval:   ap.memInterval,
	}

	// Lista arquivos de profile existentes
	if files, err := filepath.Glob(filepath.Join(ap.outputDir, "*.prof")); err == nil {
		stats.ProfileFiles = files
	}

	return stats
}

// ProfilerStats estatísticas do profiler
type ProfilerStats struct {
	Enabled          bool
	IsRunning        bool
	OutputDirectory  string
	ProfilesCreated  int
	CPUProfileActive bool
	MemoryInterval   time.Duration
	ProfileFiles     []string
}

// AutoTrigger configura triggers automáticos baseados em métricas
type AutoTrigger struct {
	profiler      *AutoProfiler
	metricsCollector *MetricsCollector
	cpuThreshold  float64
	memThreshold  float64
	enabled       bool
	checking      bool
	stopCh        chan bool
}

// NewAutoTrigger cria um novo sistema de triggers automáticos
func NewAutoTrigger(profiler *AutoProfiler, collector *MetricsCollector, cpuThreshold, memThreshold float64) *AutoTrigger {
	return &AutoTrigger{
		profiler:         profiler,
		metricsCollector: collector,
		cpuThreshold:     cpuThreshold,
		memThreshold:     memThreshold,
		enabled:          true,
		stopCh:           make(chan bool),
	}
}

// StartMonitoring inicia monitoramento para triggers automáticos
func (at *AutoTrigger) StartMonitoring() {
	if !at.enabled || at.checking {
		return
	}

	at.checking = true

	go func() {
		ticker := time.NewTicker(30 * time.Second) // Verifica a cada 30 segundos
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				at.checkThresholds()

			case <-at.stopCh:
				at.checking = false
				return
			}
		}
	}()
}

// checkThresholds verifica se os thresholds foram ultrapassados
func (at *AutoTrigger) checkThresholds() {
	metrics := at.metricsCollector.GetMetrics()

	// Verifica CPU
	if metrics.CPUUsage.Usage > at.cpuThreshold {
		at.profiler.StartCPUProfile()
		
		// Para CPU profile depois de 30 segundos
		go func() {
			time.Sleep(30 * time.Second)
			at.profiler.StopCPUProfile()
		}()
	}

	// Verifica memória
	if metrics.MemoryUsage.AllocMB > at.memThreshold {
		at.profiler.CreateMemoryProfile()
	}

	// Se ambos estão altos, cria snapshot completo
	if metrics.CPUUsage.Usage > at.cpuThreshold && metrics.MemoryUsage.AllocMB > at.memThreshold {
		at.profiler.CreateFullSnapshot()
	}
}

// Stop para o sistema de triggers
func (at *AutoTrigger) Stop() {
	if at.checking {
		close(at.stopCh)
	}
}