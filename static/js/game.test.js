'use strict';
/* Browser-game tests for static/js/game.js. Mirrors the Go suite in
 * main_test.go (full-run completion, wrong-key behavior, miss penalty)
 * against the offline JS core. Run: node --test static/js/game.test.js */
const { describe, it } = require('node:test');
const assert = require('node:assert/strict');

const g = require('./game.js');

function memStorage(initial) {
    const data = Object.assign({}, initial);
    return {
        getItem: (k) => (k in data ? data[k] : null),
        setItem: (k, v) => { data[k] = String(v); },
        removeItem: (k) => { delete data[k]; },
        _data: data,
    };
}

function throwingStorage() {
    const boom = () => { throw new Error('denied'); };
    return { getItem: boom, setItem: boom, removeItem: boom };
}

/* Deterministic rand: cycles through values. */
function seq(values) {
    let i = 0;
    return () => values[(i++) % values.length];
}

function typeText(game, text, now) {
    let t = now || 1000000;
    for (const ch of text) {
        const r = g.handleKey(game, ch, t);
        t += 100;
        assert.equal(r.ignored, false);
        assert.equal(r.correct, true, `key ${JSON.stringify(ch)} should be correct`);
    }
    return t;
}

describe('full drill completion in order (mirrors TestGameCompletesEveryLessonInOrder)', () => {
    it('completes drills and regenerates the next prompt', () => {
        const game = g.createGame({ storage: g.createStorage(memStorage()), stepIdx: 0 });
        // rand 0.99: no surprise (step 0 never surprises), pattern branch (>0.3).
        g.generateDrill(game, seq([0.99, 0.1, 0.2, 0.3]));
        assert.ok(game.currentText.length > 0);
        const first = game.currentText;
        typeText(game, first);
        // After the final char a fresh drill is generated and progress resets.
        assert.equal(game.charIndex, 0);
        assert.ok(game.currentText.length > 0);
        assert.equal(game.streak, first.length);
        assert.equal(g.accuracy(game), 100);
    });

    it('promotes through every stage with sustained accurate play', () => {
        const store = memStorage();
        const game = g.createGame({ storage: g.createStorage(store), stepIdx: 0 });
        const rand = seq([0.5]); // pattern branch, deterministic picks
        g.generateDrill(game, rand);
        let guard = 0;
        let t = 1000000;
        while (game.stepIdx < g.ADVENTURE_STEPS.length - 1 && guard++ < 20000) {
            // Type the live prompt: promotions regenerate the drill mid-loop.
            const r = g.handleKey(game, game.currentText[game.charIndex], t);
            t += 100;
            assert.equal(r.correct, true);
        }
        assert.equal(game.stepIdx, g.ADVENTURE_STEPS.length - 1);
        assert.equal(store._data[g.STEP_STORAGE_KEY], String(g.ADVENTURE_STEPS.length - 1));
    });
});

describe('wrong key behavior (mirrors TestWrongKeyDoesNotAdvanceAndNextPromptIsExactKey)', () => {
    it('does not advance, counts the mistake, keeps the same prompt', () => {
        const game = g.createGame({ storage: g.createStorage(memStorage()), stepIdx: 0 });
        game.currentText = 'asdf';
        game.charIndex = 0;
        const wrong = g.handleKey(game, 'x', 1000000);
        assert.equal(wrong.correct, false);
        assert.equal(game.charIndex, 0);
        assert.equal(game.totalHits, 1);
        assert.equal(game.correctHits, 0);
        assert.equal(game.streak, 0);
        assert.equal(game.multiplier, 1);
        assert.equal(game.currentText[game.charIndex], 'a');

        const right = g.handleKey(game, 'a', 1000100);
        assert.equal(right.correct, true);
        assert.equal(game.charIndex, 1);
        assert.equal(game.currentText[game.charIndex], 's');
    });

    it('ignores modifier keys without touching stats', () => {
        const game = g.createGame({ storage: g.createStorage(memStorage()), stepIdx: 0 });
        g.generateDrill(game, seq([0.5]));
        for (const mod of ['Shift', 'Control', 'Alt', 'Meta', 'CapsLock']) {
            const r = g.handleKey(game, mod, 1000000);
            assert.equal(r.ignored, true);
        }
        assert.equal(game.totalHits, 0);
        assert.equal(game.startTime, null);
    });
});

describe('miss penalty (mirrors TestMissRespawnsTheSamePromptAndAppliesPenalty)', () => {
    it('boss timeout resets streak and multiplier and respawns a drill', () => {
        const game = g.createGame({ storage: g.createStorage(memStorage()), stepIdx: 1 });
        g.generateDrill(game, seq([0.5]));
        game.streak = 12;
        game.multiplier = 2;
        g.bossTimeout(game);
        assert.equal(game.streak, 0);
        assert.equal(game.multiplier, 1);
        assert.equal(game.charIndex, 0);
        assert.ok(game.currentText.length > 0);
    });
});

describe('multiplier progression', () => {
    it('follows the 8/15/25 streak thresholds', () => {
        assert.equal(g.multiplierForStreak(0), 1);
        assert.equal(g.multiplierForStreak(8), 1);
        assert.equal(g.multiplierForStreak(9), 2);
        assert.equal(g.multiplierForStreak(15), 2);
        assert.equal(g.multiplierForStreak(16), 3);
        assert.equal(g.multiplierForStreak(25), 3);
        assert.equal(g.multiplierForStreak(26), 4);
    });
});

describe('promotion gate (mirrors accuracy/mastery rules)', () => {
    it('requires 25 stage hits at 95%+ accuracy', () => {
        const game = g.createGame({ storage: g.createStorage(memStorage()), stepIdx: 0 });
        game.stageTotalHits = 24;
        game.stageCorrectHits = 24;
        assert.equal(g.checkPromotion(game), false);
        game.stageTotalHits = 25;
        game.stageCorrectHits = 23; // 92%
        assert.equal(g.checkPromotion(game), false);
        game.stageCorrectHits = 24; // 96%
        assert.equal(g.checkPromotion(game), true);
        assert.equal(game.stepIdx, 1);
        assert.equal(game.stageTotalHits, 0);
    });
});

describe('stats', () => {
    it('computes WPM from correct chars over elapsed minutes', () => {
        const game = g.createGame({ storage: g.createStorage(memStorage()), stepIdx: 0 });
        assert.equal(g.wpm(game, 2000000), 0);
        game.startTime = 1000000;
        game.correctHits = 250; // 50 words in 1 minute
        assert.equal(g.wpm(game, 1060000), 50);
    });

    it('accuracy starts at 100 and rounds', () => {
        const game = g.createGame({ storage: g.createStorage(memStorage()), stepIdx: 0 });
        assert.equal(g.accuracy(game), 100);
        game.totalHits = 3;
        game.correctHits = 2;
        assert.equal(g.accuracy(game), 67);
    });
});

describe('offline storage safety', () => {
    it('sound is off by default and the preference persists', () => {
        const store = memStorage();
        const game = g.createGame({ storage: g.createStorage(store) });
        assert.equal(game.soundEnabled, false);
        g.setSoundEnabled(game, true);
        assert.equal(store._data[g.SOUND_STORAGE_KEY], '1');
        const reloaded = g.createGame({ storage: g.createStorage(store) });
        assert.equal(reloaded.soundEnabled, true);
    });

    it('never throws when storage is missing or denied', () => {
        const denied = g.createStorage(throwingStorage());
        assert.equal(denied.get('x'), null);
        denied.set('x', '1');
        denied.remove('x');
        assert.equal(g.loadStepIdx(denied), 0);
        const game = g.createGame({ storage: denied });
        assert.equal(game.stepIdx, 0);
        assert.equal(game.soundEnabled, false);
        g.setSoundEnabled(game, true);
        g.resetProgress(game);
    });

    it('clamps corrupt saved stage values', () => {
        assert.equal(g.loadStepIdx(g.createStorage(memStorage({ [g.STEP_STORAGE_KEY]: 'nope' }))), 0);
        assert.equal(g.loadStepIdx(g.createStorage(memStorage({ [g.STEP_STORAGE_KEY]: '-3' }))), 0);
        assert.equal(g.loadStepIdx(g.createStorage(memStorage({ [g.STEP_STORAGE_KEY]: '99' }))), 0);
        assert.equal(g.loadStepIdx(g.createStorage(memStorage({ [g.STEP_STORAGE_KEY]: '4' }))), 4);
    });

    it('reset clears the saved stage', () => {
        const store = memStorage({ [g.STEP_STORAGE_KEY]: '5' });
        const game = g.createGame({ storage: g.createStorage(store) });
        assert.equal(game.stepIdx, 5);
        g.resetProgress(game);
        assert.equal(game.stepIdx, 0);
        assert.ok(!(g.STEP_STORAGE_KEY in store._data));
    });
});

describe('drill generation', () => {
    it('only uses keys unlocked up to the current stage', () => {
        const game = g.createGame({ storage: g.createStorage(memStorage()), stepIdx: 0 });
        const pool = new Set(g.getActiveKeySet(game));
        for (let i = 0; i < 50; i++) {
            g.generateDrill(game, seq([0.01, 0.11, 0.23, 0.37, 0.53, 0.71, 0.83, 0.97]));
            for (const ch of game.currentText) {
                if (ch === ' ') continue;
                assert.ok(pool.has(ch), `locked key ${JSON.stringify(ch)} in drill`);
            }
        }
    });

    it('grandmaster stage draws from the full keyboard', () => {
        const game = g.createGame({ storage: g.createStorage(memStorage()), stepIdx: 11 });
        assert.ok(g.getActiveKeySet(game).length > 40);
        g.generateDrill(game, seq([0.5]));
        assert.ok(game.currentText.length > 0);
    });
});
