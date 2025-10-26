<script setup lang="ts">
import { ref } from "vue";
import Music from "./Music.vue";
import Books from "./Books.vue";
import Movies from "./Movies.vue";
import type { TabType } from "../services/types";

const activeTab = ref<TabType>("music");

const tabs = [
    { label: "音乐", value: "music" },
    { label: "影视", value: "movies" },
    { label: "书籍", value: "books" },
];

const handleTabChange = (tab: TabType) => {
    activeTab.value = tab;
};
</script>

<template>
    <div class="main-page">
        <header class="header">
            <h1 class="title">书影音</h1>
            <p class="subtitle">记录你的音乐、影视和书籍</p>
        </header>

        <main class="content">
            <div class="tabs-container">
                <div class="tabs">
                    <button
                        v-for="tab in tabs"
                        :key="tab.value"
                        :class="[
                            'tab-button',
                            { active: activeTab === tab.value },
                        ]"
                        @click="handleTabChange(tab.value as TabType)"
                    >
                        {{ tab.label }}
                    </button>
                </div>
            </div>

            <div class="tab-content">
                <Music v-if="activeTab === 'music'" />
                <Books v-if="activeTab === 'books'" />
                <Movies v-if="activeTab === 'movies'" />
            </div>
        </main>
    </div>
</template>

<style scoped>
.main-page {
    min-height: 100vh;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.header {
    text-align: center;
    padding: 40px 20px 20px;
    color: white;
}

.title {
    font-size: 2.5rem;
    font-weight: 700;
    margin: 0 0 8px 0;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.subtitle {
    font-size: 1.1rem;
    opacity: 0.9;
    margin: 0;
    font-weight: 300;
}

.content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px 40px;
}

.tabs-container {
    background: white;
    border-radius: 12px;
    margin-bottom: 24px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.tabs {
    display: flex;
    padding: 8px;
}

.tab-button {
    flex: 1;
    padding: 12px 16px;
    border: none;
    background: transparent;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 500;
    color: #666;
    cursor: pointer;
    transition: all 0.3s ease;
}

.tab-button:hover {
    background: #f5f5f5;
    color: #333;
}

.tab-button.active {
    background: #409eff;
    color: white;
    box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
}

.tab-content {
    background: white;
    border-radius: 12px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
    overflow: hidden;
}

@media (max-width: 768px) {
    .title {
        font-size: 2rem;
}
    .content {
        padding: 0 16px 20px;
    }

    .tab-button {
        padding: 10px 12px;
        font-size: 0.9rem;
    }
}

@media (max-width: 480px) {
    .header {
        padding: 30px 16px 15px;
    }

    .title {
        font-size: 1.75rem;
    }

    .subtitle {
        font-size: 1rem;
    }

    .tabs {
        padding: 6px;
    }

    .tab-button {
        padding: 8px 10px;
        font-size: 0.85rem;
    }
}

@media (min-width: 1200px) {
    .content {
        max-width: 1400px;
        padding: 0 40px 60px;
    }

    .header {
        padding: 60px 20px 30px;
    }

    .title {
        font-size: 3rem;
    }

    .subtitle {
        font-size: 1.2rem;
    }
}
</style>