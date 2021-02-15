<template>
  <header>
    <div class="main">YouTube Video Picker</div>
  </header>
  <main>
    <div :class="['url-input', isInvalid ? 'invalid' : undefined]">
      <div class="label">Playlist URL</div>
      <input type="text" placeholder="https://www.youtube.com/playlist?list=..." v-model="playlistURL">
    </div>
    <div class="time-picker">
      <div class="label">Timespans</div>
      <div class="picker">
        <div class="row" v-for="timespan in timespans" :key="timespan.display">
          <input type="checkbox" v-model="selectedTimespans" :value="timespan"><span class="time">{{ timespan.display }}</span>
        </div>
      </div>
    </div>
    <button :class="['pick', pickerDisabled ? 'disabled' : undefined, isLoading ? 'loading' : undefined]" :disabled="pickerDisabled" @click="pickVideo">
      <img src="/die.png" alt="" width="24"><span>Pick Random Video</span>
    </button>
    <div v-if="resultVideo" class="result">
      <a :href="resultVideoURL">
        <img :src="resultVideo.ThumbnailURL" :alt="resultVideo.Title">
        <div class="duration">{{ resultVideoDurationDisplay }}</div>
      </a>
      <div class="title">{{ resultVideo.Title }}</div>
    </div>
  </main>
</template>

<script lang="ts">
import { Vue } from 'vue-class-component';
import Axios from 'axios';

interface VideoData {
  Duration: number;
  ID: string;
  ThumbnailURL: string;
  Title: string
}

interface PlaylistData {
  ID: string;
  Videos: VideoData[];
  Count: number;
  CacheTime: number;
  Error?: true;
}

export default class App extends Vue {
  protected playlistData: PlaylistData | null = null;

  protected playlistURL: string = '';
  protected isInvalid: boolean = false;

  protected resultVideo: VideoData | null = null;

  protected selectedTimespans: ({ min: number, max: number, display: string })[] = [];

  protected isLoading: boolean = false;

  protected get timespans() {
    const spans = [];

    for (let i = 0; i < 6; ++i) {
      spans.push({
        min: i * 5 * 60,
        max: (i + 1) * 5 * 60,
        display: `${i * 5} - ${(i + 1) * 5} min`
      });
    }

    spans.push({
      min: 7 * 5 * 60,
      max: Number.MAX_SAFE_INTEGER,
      display: '30+ min'
    });

    return spans;
  }

  protected get pickerDisabled() {
    return this.isLoading || this.isInvalid || this.selectedTimespans.length === 0;
  }

  protected get resultVideoURL() {
    return `https://youtube.com/watch?v=${this.resultVideo?.ID}`;
  }

  protected get resultVideoDurationMinutes() {
    return Math.floor((this.resultVideo?.Duration || 0) / 60);
  }

  protected get resultVideoDurationSeconds() {
    return (this.resultVideo?.Duration || 0 - (this.resultVideoDurationMinutes * 60)) % 60
  }

  protected get resultVideoDurationDisplay() {
    return `${('' + this.resultVideoDurationMinutes).padStart(2, '0')}:${('' + this.resultVideoDurationSeconds).padStart(2, '0')}`
  }

  public mounted() {
    this.playlistURL = localStorage.getItem('ytvp') || '';

    if (this.playlistURL !== '') {
      this.fetchVideos();
    }
  }

  protected get playlistId() {
    if (/https:\/\//.test(this.playlistURL)) {
      return this.playlistURL.split('list=')[1];
    }

    return this.playlistURL;
  }

  protected pickVideo() {
    this.isLoading = true;

    localStorage.setItem('ytvp', this.playlistURL);

    this.fetchVideos().then((data) => {
      if (data) {
        const fittingVideos = data.Videos.filter((v) => {
          for (let i = 0; i < this.selectedTimespans.length; ++i) {
            if (v.Duration >= this.selectedTimespans[i].min && v.Duration <= this.selectedTimespans[i].max) {
              return true;
            }
          }

          return false;
        });

        this.resultVideo = fittingVideos[Math.floor(Math.random() * fittingVideos.length)];
        this.isLoading = false;
      }
    })
  }

  protected fetchVideos() {
    this.isInvalid = false;

    if (this.playlistData && this.playlistData.ID === this.playlistId) {
      return Promise.resolve(this.playlistData);
    }

    return Axios.get(`${process.env.VUE_APP_API_URL}/lists/${this.playlistId}/videos`, {
      validateStatus(status) {
        return status < 400;
      }
    })
      .then(response => {
        this.playlistData = {
          ID: this.playlistId,
          ...response.data
        };

        return this.playlistData
      })
      .catch(() => {
        this.isInvalid = true;
      })
  }
}
</script>

<style lang="scss">
@import url('https://fonts.googleapis.com/css2?family=Pacifico&family=Lato&display=swap');

* {
  box-sizing: border-box;
}

body {
  margin: 0;

  background-color: #fbfbfb;
  font-family: 'Lato', sans-serif;
}

#app {
  display: flex;
  align-items: center;
  flex-direction: column;
  width: 100%;

  margin-bottom: 24px;
}

header, main {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;

  width: 100%;
}

header {
  font-size: 30px;
  font-family: 'Pacifico', sans-serif;

  padding: 12px 0;
  margin-bottom: 12px;

  background-color: #284E7D;

  color: white;

  .sub {
    font-size: .35em;
    font-family: 'Lato', sans-serif;
    color: #ccc;
  }
}

main {
  max-width: 800px;
  text-align: center;
}

.label {
  margin: 0 0 8px 4px;
}

.url-input {
  padding: 0 16px;
  width: 100%;

  input {
    width: 100%;

    background: white;
    border: 1px solid #ddd;
    border-radius: 3px;

    padding: 8px;

    font-size: 20px;

    text-align: center;

    &::placeholder {
      color: rgba(0, 0, 0, .25);
    }
  }

  &.invalid {
    label {
      color: red;
    }

    input {
      border-color: red;
    }
  }
}

.time-picker {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;

  margin-top: 20px;

  width: 100%;

  .picker {
    width: auto;

    display: flex;
    flex-direction: column;

    justify-content: center;
    align-items: center;

    .row {
      width: 100%;
      text-align: left;
      font-size: 16px;
      line-height: 22px;

      .time {
        padding-left: 6px;
      }
    }
  }
}

button.pick {
  margin-top: 20px;
  padding: 10px 12px;

  background: #284E7D;

  border-radius: 3px;
  border: 1px solid #13253B;

  font-size: 20px;
  font-family: 'Lato', sans-serif;
  color: white;

  display: flex;

  &.disabled {
    opacity: .5;
  }

  &.loading {
    @keyframes spin {
      from { transform: rotate(0deg); }
      to { transform: rotate(360deg); }
    }

    img {
      animation-name: spin;
      animation-iteration-count: infinite;
      animation-duration: 1.25s;
    }
  }

  img {
    display: block;
  }

  span {
    margin-left: 8px;
  }
}

.result {
  border-radius: 3px;
  overflow: hidden;
  
  margin-top: 20px;

  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 90%;

  img {
    display: block;
    width: 100%;
  }

  a {
    display: block;
    width: 100%;

    position: relative;

    .duration {
      position: absolute;
      bottom: 8px;
      right: 8px;

      color: #f1f1f1;

      padding: 2px 4px;
      font-size: 13px;
      font-weight: bold;

      background: rgba(0, 0, 0, .8);
      border-radius: 3px;
    }
  }

  .title {
    border: 1px solid #ccc;
    border-top: 0;

    width: 100%;

    color: #111;

    background-color: #eee;
    padding: 8px;
  }
}
</style>
