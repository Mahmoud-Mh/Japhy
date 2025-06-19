<template>
  <div class="game-container">
    <div class="game-header">
      <h1 class="game-title">🐚 Jeu du Bonneteau 🐚</h1>
      <p class="game-subtitle">Trouvez la perle cachée sous l'un des coquillages !</p>
      
      <!-- Control Buttons -->
      <div class="control-buttons">
        <button @click="showInstructions = !showInstructions" class="instructions-toggle">
          {{ showInstructions ? '❌ Fermer' : '❓ Comment jouer' }}
        </button>
        <button @click="soundEnabled = !soundEnabled" class="sound-toggle" :class="{ muted: !soundEnabled }">
          {{ soundEnabled ? '🔊 Son' : '🔇 Muet' }}
        </button>
      </div>
      
      <!-- Instructions Panel -->
      <div v-if="showInstructions" class="instructions-panel">
        <p><strong>Règles du jeu :</strong></p>
        <ol>
          <li>Regardez bien sous quel coquillage se cache la perle</li>
          <li>Les coquillages sont mélangés</li>
          <li>Cliquez sur le coquillage où vous pensez que se trouve la perle</li>
          <li>Gagnez des points et battez votre record de victoires consécutives !</li>
        </ol>
      </div>
      
      <!-- Difficulty Selector -->
      <div class="difficulty-selector">
        <span class="difficulty-label">Difficulté :</span>
        <button 
          v-for="level in (['easy', 'medium', 'hard'] as const)" 
          :key="level"
          @click="difficulty = level"
          class="difficulty-button"
          :class="{ active: difficulty === level }"
        >
          {{ level === 'easy' ? '😌 Facile' : level === 'medium' ? '😐 Moyen' : '😈 Difficile' }}
        </button>
      </div>
      
      <div class="score-board">
        <div class="score-item">
          <span class="score-label">Victoires:</span>
          <span class="score-value">{{ wins }}</span>
        </div>
        <div class="score-item">
          <span class="score-label">Défaites:</span>
          <span class="score-value">{{ losses }}</span>
        </div>
        <div class="score-item">
          <span class="score-label">Taux de réussite:</span>
          <span class="score-value">{{ winRate }}%</span>
        </div>
        <div class="score-item">
          <span class="score-label">Série actuelle:</span>
          <span class="score-value streak">{{ currentStreak }}</span>
        </div>
        <div class="score-item">
          <span class="score-label">Meilleure série:</span>
          <span class="score-value best-streak">{{ bestStreak }}</span>
        </div>
      </div>
    </div>

    <div class="game-area">
      <!-- Pearl display phase -->
      <div v-if="gamePhase === 'reveal'" class="reveal-phase">
        <p class="phase-instruction">🔍 Regardez bien ! La perle est sous le coquillage {{ correctShell + 1 }}</p>
        <div class="shells-container">
          <div 
            v-for="(_, index) in shells" 
            :key="index" 
            class="shell-wrapper"
            :class="{ 'has-pearl': index === correctShell }"
          >
            <div class="shell reveal-shell">🐚</div>
            <div v-if="index === correctShell" class="pearl reveal-pearl">🫧</div>
          </div>
        </div>
        <button @click="startShuffle" class="action-button shuffle-button">
          Mélanger les coquillages
        </button>
      </div>

      <!-- Shuffling phase -->
      <div v-if="gamePhase === 'shuffle'" class="shuffle-phase">
        <p class="phase-instruction">✨ Mélange en cours... Suivez bien le coquillage avec la perle !</p>
        <div class="shells-container shuffling">
          <div 
            v-for="(_, index) in shells" 
            :key="index" 
            class="shell-wrapper shuffling"
            :style="{ 
              transform: `translateX(${(shells[index].position - index) * 120}px)`,
              zIndex: shells[index].position > index ? 10 : 1
            }"
          >
            <div class="shell shuffle-shell">🐚</div>
          </div>
        </div>
      </div>

      <!-- Guessing phase -->
      <div v-if="gamePhase === 'guess'" class="guess-phase">
        <p class="phase-instruction">🤔 Où se trouve la perle ? Cliquez sur un coquillage !</p>
        <div class="shells-container">
          <div 
            v-for="(_, index) in shells" 
            :key="index" 
            class="shell-wrapper clickable"
            @click="makeGuess(index)"
          >
            <div class="shell guess-shell">🐚</div>
          </div>
        </div>
      </div>

      <!-- Result phase -->
      <div v-if="gamePhase === 'result'" class="result-phase">
        <div class="result-message" :class="{ 'win': lastResult === 'win', 'lose': lastResult === 'lose' }">
          <h2 v-if="lastResult === 'win'">🎉 Félicitations ! Vous avez trouvé la perle ! 🎉</h2>
          <h2 v-else>😞 Dommage ! La perle était sous le coquillage {{ pearlFinalPosition + 1 }}</h2>
        </div>
        
        <div class="shells-container">
          <div 
            v-for="(_, index) in shells" 
            :key="index" 
            class="shell-wrapper"
            :class="{ 
              'has-pearl': index === pearlFinalPosition,
              'guessed': index === userGuess,
              'correct-guess': index === pearlFinalPosition && index === userGuess,
              'wrong-guess': index === userGuess && index !== pearlFinalPosition
            }"
          >
            <div class="shell result-shell">🐚</div>
            <div v-if="index === pearlFinalPosition" class="pearl result-pearl">🫧</div>
            <div v-if="index === userGuess && index !== pearlFinalPosition" class="wrong-indicator">❌</div>
            <div v-if="index === userGuess && index === pearlFinalPosition" class="correct-indicator">✅</div>
          </div>
        </div>
        
        <button @click="newGame" class="action-button new-game-button">
          Nouvelle partie
        </button>
      </div>

      <!-- Loading phase -->
      <div v-if="gamePhase === 'loading'" class="loading-phase">
        <div class="loading-spinner">⏳</div>
        <p class="phase-instruction">Préparation de la partie...</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'

// Game state
const gamePhase = ref<'loading' | 'reveal' | 'shuffle' | 'guess' | 'result'>('loading')
const correctShell = ref<number>(0)
const userGuess = ref<number>(-1)
const lastResult = ref<'win' | 'lose' | null>(null)
const wins = ref<number>(0)
const losses = ref<number>(0)
const currentStreak = ref<number>(0)
const bestStreak = ref<number>(0)
const difficulty = ref<'easy' | 'medium' | 'hard'>('medium')
const showInstructions = ref<boolean>(false)

// Sound system
const audioContext = ref<AudioContext | null>(null)
const soundEnabled = ref<boolean>(true)

const initAudioContext = () => {
  if (!audioContext.value) {
    audioContext.value = new (window.AudioContext || (window as any).webkitAudioContext)()
  }
}

const playSound = (type: 'click' | 'shuffle' | 'win' | 'lose' | 'reveal') => {
  if (!soundEnabled.value || !audioContext.value) return
  
  const ctx = audioContext.value
  const oscillator = ctx.createOscillator()
  const gainNode = ctx.createGain()
  
  oscillator.connect(gainNode)
  gainNode.connect(ctx.destination)
  
  const now = ctx.currentTime
  
  switch (type) {
    case 'click':
      // Short click sound
      oscillator.frequency.setValueAtTime(800, now)
      oscillator.frequency.exponentialRampToValueAtTime(600, now + 0.1)
      gainNode.gain.setValueAtTime(0.3, now)
      gainNode.gain.exponentialRampToValueAtTime(0.01, now + 0.1)
      oscillator.start(now)
      oscillator.stop(now + 0.1)
      break
      
    case 'shuffle':
      // Swoosh sound
      oscillator.frequency.setValueAtTime(200, now)
      oscillator.frequency.exponentialRampToValueAtTime(800, now + 0.2)
      oscillator.frequency.exponentialRampToValueAtTime(150, now + 0.4)
      gainNode.gain.setValueAtTime(0.2, now)
      gainNode.gain.exponentialRampToValueAtTime(0.01, now + 0.4)
      oscillator.start(now)
      oscillator.stop(now + 0.4)
      break
      
    case 'win':
      // Triumphant ascending sound
      const frequencies = [523, 659, 784, 1047] 
      frequencies.forEach((freq, i) => {
        const osc = ctx.createOscillator()
        const gain = ctx.createGain()
        osc.connect(gain)
        gain.connect(ctx.destination)
        osc.frequency.setValueAtTime(freq, now + i * 0.15)
        gain.gain.setValueAtTime(0.25, now + i * 0.15)
        gain.gain.exponentialRampToValueAtTime(0.01, now + i * 0.15 + 0.3)
        osc.start(now + i * 0.15)
        osc.stop(now + i * 0.15 + 0.3)
      })
      break
      
    case 'lose':
      // Descending sad sound
      oscillator.frequency.setValueAtTime(400, now)
      oscillator.frequency.exponentialRampToValueAtTime(200, now + 0.5)
      oscillator.frequency.exponentialRampToValueAtTime(150, now + 1.0)
      gainNode.gain.setValueAtTime(0.3, now)
      gainNode.gain.exponentialRampToValueAtTime(0.01, now + 1.0)
      oscillator.start(now)
      oscillator.stop(now + 1.0)
      break
      
    case 'reveal':
      // Magical reveal sound
      oscillator.frequency.setValueAtTime(440, now)
      oscillator.frequency.exponentialRampToValueAtTime(880, now + 0.3)
      oscillator.frequency.exponentialRampToValueAtTime(440, now + 0.6)
      gainNode.gain.setValueAtTime(0.2, now)
      gainNode.gain.exponentialRampToValueAtTime(0.01, now + 0.6)
      oscillator.start(now)
      oscillator.stop(now + 0.6)
      break
  }
}

// Shell positions for animation
const shells = ref([
  { id: 0, position: 0 }, // Left shell
  { id: 1, position: 1 }, // Center shell  
  { id: 2, position: 2 }  // Right shell
])

// Computed properties
const winRate = computed(() => {
  const total = wins.value + losses.value
  return total > 0 ? Math.round((wins.value / total) * 100) : 0
})

const difficultySettings = computed(() => {
  switch (difficulty.value) {
    case 'easy': return { shuffles: 6, speed: 500 }
    case 'medium': return { shuffles: 8, speed: 400 }
    case 'hard': return { shuffles: 12, speed: 300 }
    default: return { shuffles: 8, speed: 400 }
  }
})

// Computed property to find where the pearl-carrying shell ended up
const pearlFinalPosition = computed(() => {
  if (gamePhase.value === 'result') {
    // Find where the shell with the pearl ended up
    const pearlShell = shells.value.find(shell => shell.id === correctShell.value)
    return pearlShell ? pearlShell.position : -1
  }
  return correctShell.value // During reveal phase, show at original position
})

// API call to get random number
const fetchRandomNumber = async (): Promise<number> => {
  try {
    const response = await fetch('https://www.random.org/integers/?num=1&min=0&max=2&col=1&base=10&format=plain&rnd=new')
    const text = await response.text()
    return parseInt(text.trim())
  } catch (error) {
    console.error('Error fetching random number:', error)
    // Fallback to local random if API fails
    return Math.floor(Math.random() * 3)
  }
}

// Initialize new game
const newGame = async () => {
  gamePhase.value = 'loading'
  userGuess.value = -1
  lastResult.value = null
  
  // Initialize audio context on first interaction
  initAudioContext()
  
  // Reset shell positions
  shells.value = [
    { id: 0, position: 0 },
    { id: 1, position: 1 },
    { id: 2, position: 2 }
  ]
  
  // Get random shell position from API
  correctShell.value = await fetchRandomNumber()
  
  // Start reveal phase
  gamePhase.value = 'reveal'
  
  // Play reveal sound
  setTimeout(() => playSound('reveal'), 500)
}

// Start shuffling animation
const startShuffle = () => {
  gamePhase.value = 'shuffle'
  playSound('click') // Click sound when starting shuffle
  
  // Animate shell shuffling based on difficulty
  const { shuffles, speed } = difficultySettings.value
  let currentShuffle = 0
  
  const shuffleInterval = setInterval(() => {
    // Select two distinct shells
    let a = Math.floor(Math.random() * 3)
    let b = Math.floor(Math.random() * 3)
    while (a === b) {
      b = Math.floor(Math.random() * 3)
    }

    // Swap their positions
    const temp = shells.value[a].position
    shells.value[a].position = shells.value[b].position
    shells.value[b].position = temp

    // Play shuffle sound on each swap
    playSound('shuffle')

    currentShuffle++
    if (currentShuffle >= shuffles) {
      clearInterval(shuffleInterval)
      setTimeout(() => {
        gamePhase.value = 'guess'
      }, 500)
    }
  }, speed)
}

// Handle user guess
const makeGuess = (shellIndex: number) => {
  if (gamePhase.value !== 'guess') return
  
  userGuess.value = shellIndex
  playSound('click') // Click sound on shell selection
  
  // Find which shell is at the guessed position after shuffling
  let actualShellAtPosition = -1
  for (let i = 0; i < shells.value.length; i++) {
    if (shells.value[i].position === shellIndex) {
      actualShellAtPosition = i
      break
    }
  }
  
  // Check if correct
  if (actualShellAtPosition === correctShell.value) {
    lastResult.value = 'win'
    wins.value++
    currentStreak.value++
    if (currentStreak.value > bestStreak.value) {
      bestStreak.value = currentStreak.value
    }
    // Play win sound after a short delay
    setTimeout(() => playSound('win'), 300)
  } else {
    lastResult.value = 'lose'
    losses.value++
    currentStreak.value = 0
    // Play lose sound after a short delay
    setTimeout(() => playSound('lose'), 300)
  }
  
  gamePhase.value = 'result'
}

// Initialize game on mount
onMounted(() => {
  newGame()
})
</script>

<style scoped>
.game-container {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
  padding: 2rem;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  max-width: 800px;
  width: 90%;
  text-align: center;
}

.game-header {
  margin-bottom: 2rem;
}

.game-title {
  font-size: 2.5rem;
  color: #2c3e50;
  margin-bottom: 0.5rem;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
}

.game-subtitle {
  font-size: 1.2rem;
  color: #7f8c8d;
  margin-bottom: 1.5rem;
}

.control-buttons {
  display: flex;
  gap: 1rem;
  justify-content: center;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}

.instructions-toggle,
.sound-toggle {
  background: linear-gradient(135deg, #17a2b8 0%, #138496 100%);
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.3s ease;
  min-width: 120px;
}

.instructions-toggle:hover,
.sound-toggle:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(23, 162, 184, 0.3);
}

.sound-toggle.muted {
  background: linear-gradient(135deg, #6c757d 0%, #495057 100%);
  opacity: 0.7;
}

.sound-toggle.muted:hover {
  box-shadow: 0 4px 12px rgba(108, 117, 125, 0.3);
}

.instructions-panel {
  background: linear-gradient(135deg, #e7f3ff 0%, #f0f8ff 100%);
  border: 2px solid #b8daff;
  border-radius: 10px;
  padding: 1rem;
  margin-bottom: 1rem;
  text-align: left;
  font-size: 0.9rem;
}

.instructions-panel ol {
  margin: 0.5rem 0 0 1.5rem;
  padding: 0;
}

.instructions-panel li {
  margin-bottom: 0.3rem;
}

.difficulty-selector {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}

.difficulty-label {
  font-weight: 600;
  color: #2c3e50;
  margin-right: 0.5rem;
}

.difficulty-button {
  background: #f8f9fa;
  border: 2px solid #dee2e6;
  padding: 0.4rem 0.8rem;
  border-radius: 20px;
  cursor: pointer;
  font-size: 0.8rem;
  transition: all 0.3s ease;
}

.difficulty-button:hover {
  background: #e9ecef;
  transform: translateY(-1px);
}

.difficulty-button.active {
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  color: white;
  border-color: #28a745;
}

.score-board {
  display: flex;
  justify-content: center;
  gap: 1.5rem;
  background: linear-gradient(135deg, #f6f9fc 0%, #e9ecef 100%);
  padding: 1rem;
  border-radius: 15px;
  border: 2px solid #dee2e6;
  flex-wrap: wrap;
}

.score-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 80px;
}

.score-label {
  font-size: 0.8rem;
  color: #6c757d;
  font-weight: 600;
  text-align: center;
}

.score-value {
  font-size: 1.3rem;
  font-weight: bold;
  color: #2c3e50;
}

.score-value.streak {
  color: #fd7e14;
}

.score-value.best-streak {
  color: #dc3545;
  position: relative;
}

.score-value.best-streak::after {
  content: '🔥';
  position: absolute;
  top: -5px;
  right: -15px;
  font-size: 0.8rem;
}

.game-area {
  min-height: 300px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.phase-instruction {
  font-size: 1.3rem;
  color: #2c3e50;
  margin-bottom: 2rem;
  font-weight: 600;
}

.shells-container {
  display: flex;
  justify-content: center;
  gap: 2rem;
  margin-bottom: 2rem;
  position: relative;
  height: 120px;
  align-items: center;
}

.shell-wrapper {
  position: relative;
  transition: all 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

.shell-wrapper.shuffling {
  transition: transform 0.4s ease-in-out, z-index 0s;
}

.shell-wrapper.shuffling .shell {
  filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.3));
}

.shell-wrapper.shuffling[style*="z-index: 10"] .shell {
  transform: scale(1.05);
  filter: drop-shadow(0 6px 12px rgba(0, 0, 0, 0.4));
}

.shell-wrapper.clickable {
  cursor: pointer;
  min-height: 80px;
  min-width: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.3s ease;
}

.shell-wrapper.clickable:hover .shell,
.shell-wrapper.clickable:active .shell {
  transform: translateY(-10px) scale(1.1);
  filter: brightness(1.2);
}

.shell-wrapper.clickable:active {
  background: rgba(102, 126, 234, 0.1);
}

.shell {
  font-size: 4rem;
  transition: all 0.3s ease;
  user-select: none;
  position: relative;
  z-index: 2;
}

.pearl {
  position: absolute;
  top: 95%;
  left: 50%;
  transform: translateX(-50%);
  font-size: 1.6rem;
  z-index: 1;
  opacity: 0.9;
}

.reveal-pearl {
  animation: pearlGlow 2s ease-in-out infinite;
}

.result-pearl {
  animation: pearlPulse 1s ease-in-out infinite;
}

.wrong-indicator, .correct-indicator {
  position: absolute;
  top: -35px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 1.8rem;
  z-index: 3;
}

.correct-indicator {
  animation: successPulse 1s ease-in-out infinite;
}

.wrong-indicator {
  animation: errorShake 0.5s ease-in-out;
}

.shell-wrapper.has-pearl .shell {
  animation: haspearlGlow 2s ease-in-out infinite;
}

.shell-wrapper.correct-guess {
  animation: correctGuess 1s ease-in-out;
}

.shell-wrapper.wrong-guess {
  animation: wrongGuess 0.5s ease-in-out;
}

.action-button {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 1rem 2rem;
  font-size: 1.1rem;
  border-radius: 50px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 600;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.action-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
}

.action-button:active {
  transform: translateY(0);
}

.result-message {
  margin-bottom: 2rem;
  padding: 1.5rem;
  border-radius: 15px;
  font-weight: bold;
}

.result-message.win {
  background: linear-gradient(135deg, #d4edda 0%, #c3e6cb 100%);
  color: #155724;
  border: 2px solid #c3e6cb;
}

.result-message.lose {
  background: linear-gradient(135deg, #f8d7da 0%, #f5c6cb 100%);
  color: #721c24;
  border: 2px solid #f5c6cb;
}

.loading-phase {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.loading-spinner {
  font-size: 3rem;
  animation: spin 1s linear infinite;
}

/* Animations */
@keyframes pearlGlow {
  0%, 100% { transform: translateX(-50%) scale(1); filter: brightness(1); }
  50% { transform: translateX(-50%) scale(1.2); filter: brightness(1.5); }
}

@keyframes pearlPulse {
  0%, 100% { transform: translateX(-50%) scale(1); }
  50% { transform: translateX(-50%) scale(1.3); }
}

@keyframes haspearlGlow {
  0%, 100% { filter: brightness(1); }
  50% { filter: brightness(1.3) drop-shadow(0 0 10px gold); }
}

@keyframes correctGuess {
  0% { transform: scale(1); }
  50% { transform: scale(1.2); }
  100% { transform: scale(1); }
}

@keyframes wrongGuess {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-10px); }
  75% { transform: translateX(10px); }
}

@keyframes successPulse {
  0%, 100% { transform: translateX(-50%) scale(1); }
  50% { transform: translateX(-50%) scale(1.3); }
}

@keyframes errorShake {
  0%, 100% { transform: translateX(-50%) rotate(0deg); }
  25% { transform: translateX(-50%) rotate(-10deg); }
  75% { transform: translateX(-50%) rotate(10deg); }
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Responsive design */
@media (max-width: 768px) {
  .game-title {
    font-size: 2rem;
  }
  
  .game-subtitle {
    font-size: 1rem;
  }
  
  .shells-container {
    gap: 1rem;
  }
  
  .shell {
    font-size: 3rem;
  }
  
  .score-board {
    gap: 1rem;
    padding: 0.8rem;
  }
  
  .score-item {
    min-width: 70px;
  }
  
  .score-label {
    font-size: 0.7rem;
  }
  
  .score-value {
    font-size: 1.1rem;
  }
  
  .game-container {
    padding: 1rem;
    width: 95%;
  }
  
  .difficulty-selector {
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .difficulty-button {
    font-size: 0.9rem;
    padding: 0.5rem 1rem;
  }
  
  .instructions-panel {
    font-size: 0.8rem;
  }
  
  .phase-instruction {
    font-size: 1.1rem;
  }
}

@media (max-width: 480px) {
  .shells-container {
    gap: 0.5rem;
  }
  
  .shell {
    font-size: 2.5rem;
  }
  
  .shell-wrapper.clickable {
    min-height: 70px;
    min-width: 70px;
  }
  
  .score-board {
    gap: 0.5rem;
  }
  
  .score-item {
    min-width: 60px;
  }
}
</style> 