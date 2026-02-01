# 🔧 TASK-004 : Refactor SubscriptionService
## Implementation Guide (5 SP, J3-J4)

**Status** : READY FOR DEV  
**Assigned to** : Amelia (@bmad-dev)  
**Depends on** : TASK-001 (Clean Architecture)  
**Blocks** : TASK-002 (Tests)  
**Story Points** : 5 SP  
**Estimated timeline** : 2-3 days

---

## 📋 ACCEPTANCE CRITERIA (100% requis)

```
✅ [ ] SchedulerService gère `nextCheckAt` logic
✅ [ ] EpisodeResolver gère fetching episodes anime-sama
✅ [ ] Tests passent après refactor
✅ [ ] Aucune régression fonctionnelle
```

---

## 🎯 OBJECTIF GÉNÉRAL

Refactoriser `subscriptions.go` pour **respecter Single Responsibility Principle** :

**Avant** : `SubscriptionService` = 4 responsabilités
```
1. CRUD subscriptions
2. Fetching episodes anime-sama
3. Enqueue logic (créer jobs)
4. Sync orchestration + NextCheckAt calculation
```

**Après** : Responsabilités séparées
```
SubscriptionService         → CRUD uniquement
EpisodeResolver             → Fetch + parse episodes
EnqueueCoordinator (bonus)  → Créer jobs
SubscriptionScheduler       → Calcul NextCheckAt
```

---

## 📂 STRUCTURE FICHIERS À CRÉER/MODIFIER

### Fichiers à CRÉER

```
internal/ports/
├── job_creator.go                 ← NOUVEAU

internal/app/
├── episode_resolver.go             ← NOUVEAU
├── subscriptions.go                ← MODIFIER (réduire scope)
└── subscription_scheduler.go       ← MODIFIER (ajouter ComputeNextCheck)
```

### Fichiers à MODIFIER

```
internal/app/
├── subscriptions.go                ← Enlever logique, garder CRUD
└── subscription_scheduler.go       ← Ajouter ComputeNextCheck()
```

### Tests à CRÉER (bonus, part de TASK-002)

```
internal/app/
├── episode_resolver_test.go        ← Tests pour resolver
└── subscription_scheduler_test.go  ← Tests pour scheduler
```

---

## 🔴 STEP 1 : Créer interface `ports/JobCreator`

### Fichier à créer : `internal/ports/job_creator.go`

```go
package ports

import (
	"context"
	"encoding/json"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

// JobCreator defines the interface for creating download jobs.
type JobCreator interface {
	Create(ctx context.Context, req JobCreationRequest) (domain.Job, error)
}

// JobCreationRequest is the request structure for creating jobs.
type JobCreationRequest struct {
	Type   string          `json:"type"`
	Params json.RawMessage `json:"params,omitempty"`
}
```

### Pourquoi ?
- ✅ Abstraction pour `JobService`
- ✅ Mockable dans les tests
- ✅ Respecte Clean Architecture (Ports Layer)
- ✅ Découple SubscriptionService de JobService

### Impact
- **Fichiers affectés** : Aucun pour le moment (sera utilisé en Step 3)
- **Tests** : Pourront utiliser mockJobCreator
- **Effort** : 0.5 SP

---

## 🟠 STEP 2 : Créer `EpisodeResolver`

### Fichier à créer : `internal/app/episode_resolver.go`

**Responsabilité** : Résoudre les épisodes disponibles sur anime-sama

```go
package app

import (
	"context"
	"fmt"
	"strings"
)

// EpisodeResolver handles fetching and parsing anime episodes from anime-sama.
type EpisodeResolver struct {
	// Optionnel : injecter HTTP client si besoin custom (ex: timeouts)
}

// NewEpisodeResolver creates a new episode resolver.
func NewEpisodeResolver() *EpisodeResolver {
	return &EpisodeResolver{}
}

// ResolvedEpisodes represents the result of episode resolution.
type ResolvedEpisodes struct {
	SelectedPlayer string
	URLs           []string
	MaxEpisode     int
}

// Resolve fetches episodes from anime-sama baseURL and resolves player selection.
//
// Flow:
// 1. Fetch episodes.js from baseURL
// 2. Parse players map
// 3. Select best player (fallback if preferred unavailable)
// 4. Return URLs + max episode count
func (r *EpisodeResolver) Resolve(ctx context.Context, baseURL, preferredPlayer string) (ResolvedEpisodes, error) {
	// Step 1: Fetch episodes.js
	jsText, err := FetchEpisodesJS(ctx, baseURL)
	if err != nil {
		return ResolvedEpisodes{}, fmt.Errorf("fetch episodes.js: %w", err)
	}

	// Step 2: Parse episodes
	eps, err := ParseEpisodesJS(jsText)
	if err != nil {
		return ResolvedEpisodes{}, fmt.Errorf("parse episodes.js: %w", err)
	}

	// Step 3: Select player
	selected, urls := selectPlayer(preferredPlayer, eps.Players)

	// Step 4: Calculate max episode
	maxAvail := MaxAvailableEpisode(urls)
	if maxAvail < 0 {
		maxAvail = 0
	}

	return ResolvedEpisodes{
		SelectedPlayer: selected,
		URLs:           urls,
		MaxEpisode:     maxAvail,
	}, nil
}
```

### Onde extraire code FROM

**Atual em `SubscriptionService.Episodes()` (linhas 205-240)**:
```go
jsText, err := FetchEpisodesJS(ctx, sub.BaseURL)
if err != nil {
    return EpisodesResponse{}, err
}
eps, err := ParseEpisodesJS(jsText)
if err != nil {
    return EpisodesResponse{}, err
}

selected := sub.Player
if selected == "" || strings.EqualFold(selected, "auto") {
    selected = BestPlayer(eps.Players)
    if selected == "auto" {
        selected = ""
    }
}
urls := eps.Players[selected]
if len(urls) == 0 {
    selected = BestPlayer(eps.Players)
    urls = eps.Players[selected]
}

maxAvail := MaxAvailableEpisode(urls)
if maxAvail < 0 {
    maxAvail = 0
}
```

**→ Mover para `EpisodeResolver.Resolve()`**

### Também em `SubscriptionService.SyncOnce()` (linhas 352-390)

Même logique → utilise `EpisodeResolver.Resolve()` après

### Métodos auxiliares (já existem, podem ficar globais)

```go
// Helper functions (já em app/animesama.go ou similar)
func selectPlayer(preferred string, players map[string][]string) (string, []string)
func BestPlayer(players map[string][]string) string
func MaxAvailableEpisode(urls []string) int
// ^^ Deixar como estão, apenas refactorizar SubscriptionService para usar EpisodeResolver
```

### Effort
- **Extracting** : 1.5 SP
- **Testing** : Part of TASK-002

---

## 🟡 STEP 3 : Modifier `SubscriptionService`

### Fichier : `internal/app/subscriptions.go`

#### 3.1 Atualize struct + constructor

**ANTES** :
```go
type SubscriptionService struct {
	repo ports.SubscriptionRepository
	jobs *JobService              // ❌ Concrète
	bus  ports.EventBus
}

func NewSubscriptionService(repo ports.SubscriptionRepository, jobs *JobService, bus ports.EventBus) *SubscriptionService {
	return &SubscriptionService{repo: repo, jobs: jobs, bus: bus}
}
```

**DEPOIS** :
```go
type SubscriptionService struct {
	repo             ports.SubscriptionRepository
	jobCreator       ports.JobCreator             // ✅ Interface
	episodeResolver  *EpisodeResolver             // ✅ Novo
	bus              ports.EventBus
}

func NewSubscriptionService(
	repo ports.SubscriptionRepository,
	jobCreator ports.JobCreator,      // ← NOVO parâmetro
	episodeResolver *EpisodeResolver, // ← NOVO parâmetro
	bus ports.EventBus,
) *SubscriptionService {
	return &SubscriptionService{
		repo:            repo,
		jobCreator:      jobCreator,
		episodeResolver: episodeResolver,
		bus:             bus,
	}
}
```

#### 3.2 Atualizar `Episodes()`

**ANTES** :
```go
func (s *SubscriptionService) Episodes(ctx context.Context, id string) (EpisodesResponse, error) {
	sub, err := s.repo.Get(ctx, id)
	if err != nil {
		return EpisodesResponse{}, err
	}

	jsText, err := FetchEpisodesJS(ctx, sub.BaseURL)
	if err != nil {
		return EpisodesResponse{}, err
	}
	eps, err := ParseEpisodesJS(jsText)
	if err != nil {
		return EpisodesResponse{}, err
	}

	selected, urls := selectPlayer(sub.Player, eps.Players)
	maxAvail := MaxAvailableEpisode(urls)
	// ...
}
```

**DEPOIS** :
```go
func (s *SubscriptionService) Episodes(ctx context.Context, id string) (EpisodesResponse, error) {
	sub, err := s.repo.Get(ctx, id)
	if err != nil {
		return EpisodesResponse{}, err
	}

	// ✅ Usar EpisodeResolver
	resolved, err := s.episodeResolver.Resolve(ctx, sub.BaseURL, sub.Player)
	if err != nil {
		return EpisodesResponse{}, err
	}

	out := make([]EpisodeStatus, 0, resolved.MaxEpisode)
	for ep := 1; ep <= resolved.MaxEpisode; ep++ {
		available := false
		if ep-1 >= 0 && ep-1 < len(resolved.URLs) {
			available = strings.TrimSpace(resolved.URLs[ep-1]) != ""
		}
		out = append(out, EpisodeStatus{
			Episode:    ep,
			Available:  available,
			Scheduled:  sub.LastScheduledEpisode >= ep,
			Downloaded: sub.LastDownloadedEpisode >= ep,
		})
	}

	return EpisodesResponse{
		Subscription:        toSubscriptionDTO(sub),
		SelectedPlayer:      resolved.SelectedPlayer,
		MaxAvailableEpisode: resolved.MaxEpisode,
		Episodes:            out,
	}, nil
}
```

#### 3.3 Atualizar `EnqueueEpisodes()`

**ANTES** :
```go
func (s *SubscriptionService) EnqueueEpisodes(ctx context.Context, id string, episodes []int) (EnqueueEpisodesResponse, error) {
	if s.jobs == nil {                              // ❌ Checa JobService
		return EnqueueEpisodesResponse{}, errors.New("job service not configured")
	}
	// ...
	b, _ := json.Marshal(params)
	created, err := s.jobs.Create(ctx, CreateJobRequest{Type: "download", Params: b})
	// ...
}
```

**DEPOIS** :
```go
func (s *SubscriptionService) EnqueueEpisodes(ctx context.Context, id string, episodes []int) (EnqueueEpisodesResponse, error) {
	if s.jobCreator == nil {                        // ✅ Checa interface
		return EnqueueEpisodesResponse{}, errors.New("job creator not configured")
	}
	
	// ... resto do código igual

	b, _ := json.Marshal(params)
	created, err := s.jobCreator.Create(ctx, ports.JobCreationRequest{  // ✅ Use interface
		Type:   "download",
		Params: b,
	})
	// ...
}
```

#### 3.4 Atualizar `SyncOnce()`

**ANTES** :
```go
func (s *SubscriptionService) SyncOnce(ctx context.Context, id string, enqueue bool) (SyncResult, error) {
	// ...
	jsText, err := FetchEpisodesJS(ctx, sub.BaseURL)
	if err != nil {
		// ...
		return SyncResult{}, err
	}
	eps, err := ParseEpisodesJS(jsText)
	if err != nil {
		return SyncResult{}, err
	}

	selected := sub.Player
	if selected == "" || strings.EqualFold(selected, "auto") {
		selected = BestPlayer(eps.Players)
		// ...
	}
	urls := eps.Players[selected]
	// ... calcul maxAvail
	
	// ❌ Calcul NextCheckAt aqui
	if sub.LastScheduledEpisode < maxAvail {
		sub.NextCheckAt = now.Add(10 * time.Minute)
	} else {
		sub.NextCheckAt = now.Add(2 * time.Hour)
	}
}
```

**DEPOIS** :
```go
func (s *SubscriptionService) SyncOnce(ctx context.Context, id string, enqueue bool) (SyncResult, error) {
	sub, err := s.repo.Get(ctx, id)
	if err != nil {
		return SyncResult{}, err
	}

	// ✅ Usar EpisodeResolver
	resolved, err := s.episodeResolver.Resolve(ctx, sub.BaseURL, sub.Player)
	if err != nil {
		sub.LastCheckedAt = time.Now().UTC()
		sub.NextCheckAt = time.Now().UTC().Add(30 * time.Minute)
		sub.UpdatedAt = time.Now().UTC()
		_, _ = s.repo.Update(ctx, sub)
		return SyncResult{}, err
	}

	now := time.Now().UTC()
	sub.LastAvailableEpisode = resolved.MaxEpisode
	sub.LastCheckedAt = now

	// ✅ Delegar cálculo de NextCheckAt a SchedulerService (via méthode)
	// (implementado em STEP 4)
	// sub.NextCheckAt = s.scheduler.ComputeNextCheck(sub, resolved.MaxEpisode)

	// OU version simple inline (TBD após STEP 4)
	if sub.LastScheduledEpisode < resolved.MaxEpisode {
		sub.NextCheckAt = now.Add(10 * time.Minute)
	} else {
		sub.NextCheckAt = now.Add(2 * time.Hour)
	}

	// ... resto : enqueue logic
}
```

### Effort
- **Updating** : 1 SP
- **Testing regressions** : Part of TASK-002

---

## 🟢 STEP 4 : Ajouter `ComputeNextCheck()` à `SubscriptionScheduler`

### Fichier : `internal/app/subscription_scheduler.go`

#### 4.1 Ajouter méthode

**AVANT** :
```go
type SubscriptionScheduler struct {
	logger zerolog.Logger
	subs   *SubscriptionService
	repo   ports.SubscriptionRepository

	TickInterval time.Duration
	BatchSize    int
	Enqueue      bool
}

func (sch *SubscriptionScheduler) Run(ctx context.Context) {
	// ...
}

func (sch *SubscriptionScheduler) tick(ctx context.Context) {
	// ...
}
```

**APRÈS** :
```go
type SubscriptionScheduler struct {
	logger zerolog.Logger
	subs   *SubscriptionService
	repo   ports.SubscriptionRepository

	TickInterval time.Duration
	BatchSize    int
	Enqueue      bool
}

// ComputeNextCheck calculates when to check a subscription for new episodes.
//
// Logic:
// - If new episodes available beyond last scheduled: check in 10 minutes
// - Otherwise: check in 2 hours
func (sch *SubscriptionScheduler) ComputeNextCheck(sub domain.Subscription, maxAvailable int) time.Time {
	now := time.Now().UTC()

	if sub.LastScheduledEpisode < maxAvailable {
		// Nouveaux épisodes disponibles → check fréquent
		return now.Add(10 * time.Minute)
	}

	// À jour → check espacé
	return now.Add(2 * time.Hour)
}

func (sch *SubscriptionScheduler) Run(ctx context.Context) {
	// ... (unchanged)
}

func (sch *SubscriptionScheduler) tick(ctx context.Context) {
	// ... (unchanged, mais peut utiliser ComputeNextCheck si refactoring plus tard)
}
```

#### 4.2 (Optionnel) Utiliser dans `SyncOnce()`

Si vous voulez centraliser complètement, modifier `SyncOnce()` :

```go
// Dans SubscriptionService.SyncOnce()
// Au lieu de :
if sub.LastScheduledEpisode < resolved.MaxEpisode {
    sub.NextCheckAt = now.Add(10 * time.Minute)
} else {
    sub.NextCheckAt = now.Add(2 * time.Hour)
}

// Faire :
sub.NextCheckAt = s.scheduler.ComputeNextCheck(sub, resolved.MaxEpisode)
```

**Mais attention** : Cela nécessite injecter `*SubscriptionScheduler` dans `SubscriptionService`, créant une relation circulaire (Service ↔ Scheduler).

**Recommandation** : Laisser la logique inline dans `SyncOnce()` pour moment, mais `ComputeNextCheck()` existe comme utility pour tests + future refactor.

### Effort
- **Adding method** : 0.5 SP

---

## 🔄 STEP 5 : Vérifier tous les appels et intégration

### Files qui instantient `SubscriptionService`

**À trouver et modifier** :

Utiliser :
```bash
grep -r "NewSubscriptionService" --include="*.go"
```

**Généralement dans** :
- `cmd/asd-server/main.go` (initialization)
- `internal/adapters/httpapi/` (handlers)
- Tests (fixtures)

**Exemple - Before** :
```go
subSvc := app.NewSubscriptionService(subRepo, jobSvc, eventBus)
```

**Example - After** :
```go
episodeResolver := app.NewEpisodeResolver()
subSvc := app.NewSubscriptionService(
    subRepo,
    jobSvc,  // ← jobSvc implements ports.JobCreator, donc compatible ✅
    episodeResolver,
    eventBus,
)
```

### Effort
- **Searching + updating callsites** : 1 SP

---

## 🧪 STEP 6 : Vérifier pas de régression

### Tests existants

**Run** :
```bash
cd /home/guilhem/Anime-Sama-Downloader
go test ./internal/app -v
```

**Expected** :
- ✅ Test `TestSubscriptionService_Create_AutoLabelFromBaseURL` passe
- ✅ Aucune autre régression

### Vérifier compile

```bash
go build ./cmd/asd-server
```

**Expected** : ✅ Pas de compilation errors

### Effort
- **Testing** : 1 SP

---

## 📊 EFFORT BREAKDOWN

| Step | Effort | Description |
|------|--------|-------------|
| 1 | 0.5 SP | Créer `ports/JobCreator` |
| 2 | 1.5 SP | Créer `EpisodeResolver` |
| 3 | 1.0 SP | Modifier `SubscriptionService` |
| 4 | 0.5 SP | Ajouter `ComputeNextCheck` |
| 5 | 1.0 SP | Intégration + callsites |
| 6 | 1.0 SP | Testing + vérifications |
| **TOTAL** | **5.5 SP** | **Vs 5 SP budgeted** ⚠️ |

**Note** : 5.5 SP vs 5 SP estimé = 10% overrun possible. Risk: LOW (tasks bien définies, pas de surprises).

---

## ✅ ACCEPTANCE CRITERIA CHECKLIST

```
BEFORE merging, verify:

[ ] SchedulerService gère nextCheckAt logic
    └─ ComputeNextCheck() méthode existe + peut être utilisée
    
[ ] EpisodeResolver gère fetching episodes anime-sama
    └─ EpisodeResolver.Resolve() extrait logique
    └─ FetchEpisodesJS + ParseEpisodesJS utilisés
    
[ ] Tests passent après refactor
    └─ go test ./internal/app -v : PASS
    └─ go build ./cmd/asd-server : OK
    
[ ] Aucune régression fonctionnelle
    └─ TestSubscriptionService_Create_AutoLabelFromBaseURL : PASS
    └─ Episodes() tests existants : PASS
    └─ EnqueueEpisodes() comportement identique
    └─ SyncOnce() comportement identique
```

---

## 🚨 GOTCHAS & ALERTES

### ⚠️ Alerte 1 : Dépendance circulaire

**Si vous essayez** :
```go
type SubscriptionService struct {
    scheduler *SubscriptionScheduler  // ❌ Circulaire !
}
```

**Solution** : Laisser logique NextCheckAt inline dans SyncOnce(), ComputeNextCheck() est utility uniquement.

---

### ⚠️ Alerte 2 : Paramètre JobService → JobCreator

**Attention** : `JobService` implémente maintenant `ports.JobCreator`

**Avant d'injecter** :
```go
// Vérifier que JobService a:
func (j *JobService) Create(ctx context.Context, req JobCreationRequest) (domain.Job, error)
```

**Sinon** : Créer adapter ou renommer méthode.

---

### ⚠️ Alerte 3 : Ordre des paramètres

**Nouvelle signature** :
```go
func NewSubscriptionService(
    repo ports.SubscriptionRepository,
    jobCreator ports.JobCreator,
    episodeResolver *EpisodeResolver,
    bus ports.EventBus,
)
```

**Ancien ordre** :
```go
func NewSubscriptionService(
    repo ports.SubscriptionRepository,
    jobs *JobService,
    bus ports.EventBus,
)
```

**Attention** : Tous les appels doivent être mis à jour !

---

## 📞 REVIEW POINTS (Before merge)

**Amelia** : Quand prêt, faire PR avec :

1. ✅ Code changes bien documentés (commit messages clairs)
2. ✅ Tests existants PASSENT
3. ✅ Pas de build errors
4. ✅ PR description = checklist AC
5. ✅ Assign à Winston (architect review) AVANT merge

---

## 🎯 SUCCESS CRITERIA

✅ **TASK-004 DONE** quand :

```go
// Ces imports existent
import "github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/app"
import "github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"

// Ces types existent et fonctionnent
resolver := app.NewEpisodeResolver()
resolved, err := resolver.Resolve(ctx, baseURL, player)

svc := app.NewSubscriptionService(repo, jobCreator, resolver, bus)
episodes, err := svc.Episodes(ctx, id)

// Tests passent
// go test ./internal/app -v ✅
```

---

**Bonne chance Amelia ! 🚀**

Pour toute question technique, escalade à Winston (architect).
