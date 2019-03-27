<template>
  <header>
    <div>
      <button @click="openSidebar" :aria-label="$t('buttons.toggleSidebar')" :title="$t('buttons.toggleSidebar')" class="action">
        <i class="material-icons">menu</i>
      </button>
      <span style="font-weight: 700">File Man</span>
      <!--<img :src="logoURL" alt="File Browser">-->
      <!--<search v-if="isLogged"></search>-->
    </div>
    <div>
      <div v-if="isLogged">
        <button @click="openSearch" :aria-label="$t('buttons.search')" :title="$t('buttons.search')" class="search-button action">
          <i class="material-icons">search</i>
        </button>

        <button v-show="showSaveButton" :aria-label="$t('buttons.save')" :title="$t('buttons.save')" class="action" id="save-button">
          <i class="material-icons">save</i>
        </button>

        <button @click="openMore" id="more" :aria-label="$t('buttons.more')" :title="$t('buttons.more')" class="action">
          <i class="material-icons">more_vert</i>
        </button>
        <!-- Menu that shows on listing AND mobile when there are files selected -->
        <div id="file-selection" v-if="isMobile && isListing">
          <span v-if="selectedCount > 0">{{ selectedCount }} selected</span>
         <!-- <share-button v-show="showShareButton"></share-button>-->
          <rename-button v-show="showRenameButton"></rename-button>
         <!-- <copy-button v-show="showCopyButton"></copy-button>
          <move-button v-show="showMoveButton"></move-button>-->
          <delete-button v-show="showDeleteButton"></delete-button>
        </div>
        <!-- This buttons are shown on a dropdown on mobile phones -->
        <div id="dropdown" :class="{ active: showMore }">
          <div v-if="!isListing || !isMobile">
          <!--  <share-button v-show="showShareButton"></share-button>-->
            <rename-button v-show="showRenameButton"></rename-button>
        <!--    <copy-button v-show="showCopyButton"></copy-button>
            <move-button v-show="showMoveButton"></move-button>-->
            <delete-button v-show="showDeleteButton"></delete-button>
          </div>

       <!--   <switch-button v-show="isListing"></switch-button>-->
        <!--  <download-button v-show="showDownloadButton"></download-button>-->
          <upload-button v-show="showUpload"></upload-button>
          <info-button v-show="isFiles"></info-button>

          <!--<button v-show="isListing" @click="openSelect" :aria-label="$t('buttons.selectMultiple')" :title="$t('buttons.selectMultiple')" class="action">
            <i class="material-icons">check_circle</i>
            <span>{{ $t('buttons.select') }}</span>
          </button>-->
        </div>
        <div class="menu-account" @click="menu =!menu">
          <span class="avatar" >{{user.name}}</span>
        </div>

      </div>


      <div v-show="showOverlay" @click="resetPrompts" class="overlay"></div>
    </div>
  </header>
</template>

<script>
import Search from './Search'
import InfoButton from './buttons/Info'
import DeleteButton from './buttons/Delete'
import RenameButton from './buttons/Rename'
import UploadButton from './buttons/Upload'
import SwitchButton from './buttons/SwitchView'
import {mapGetters, mapState} from 'vuex'
import { logoURL } from '@/utils/constants'
import * as api from '@/api'
import buttons from '@/utils/buttons'
import * as auth from '@/utils/auth'
export default {
  name: 'header-layout',
  components: {
    InfoButton,
    DeleteButton,
    RenameButton,
    UploadButton,
    SwitchButton,
  },
  data: function () {
    return {
      width: window.innerWidth,
      pluginData: {
        api,
        buttons,
        'store': this.$store,
        'router': this.$router
      },
      menu: false
    }
  },
  created () {
    window.addEventListener('resize', () => {
      this.width = window.innerWidth
    })
  },
  computed: {
    ...mapGetters([
      'selectedCount',
      'isFiles',
      'isEditor',
      'isListing',
      'isLogged'
    ]),
    ...mapState([
      'req',
      'user',
      'loading',
      'reload',
      'multiple'
    ]),
    logoURL: () => logoURL,
    isMobile () {
      return this.width <= 736
    },
    showUpload () {
      return this.isListing && this.user.perm.create
    },
    showSaveButton () {
      return this.isEditor && this.user.perm.modify
    },
    showDownloadButton () {
      return this.isFiles && this.user.perm.download
    },
    showDeleteButton () {
      return this.isFiles && (this.isListing
        ? (this.selectedCount !== 0 && this.user.perm.delete)
        : this.user.perm.delete)
    },
    showRenameButton () {
      return this.isFiles && (this.isListing
        ? (this.selectedCount === 1 && this.user.perm.rename)
        : this.user.perm.rename)
    },
    showShareButton () {
      return this.isFiles && (this.isListing
        ? (this.selectedCount === 1 && this.user.perm.share)
        : this.user.perm.share)
    },
    showMoveButton () {
      return this.isFiles && (this.isListing
        ? (this.selectedCount > 0 && this.user.perm.rename)
        : this.user.perm.rename)
    },
    showCopyButton () {
      return this.isFiles && (this.isListing
        ? (this.selectedCount > 0 && this.user.perm.create)
        : this.user.perm.create)
    },
    showMore () {
      return this.isFiles && this.$store.state.show === 'more'
    },
    showOverlay () {
      return this.showMore
    }
  },
  methods: {
    openSidebar () {
      this.$store.commit('showHover', 'sidebar')
    },
    openMore () {
      this.$store.commit('showHover', 'more')
    },
    openSearch () {
      this.$store.commit('showHover', 'search')
    },
    openSelect () {
      this.$store.commit('multiple', true)
      this.resetPrompts()
    },
    resetPrompts () {
      this.$store.commit('closeHovers')
    },
      logout: auth.logout


  }
}
</script>
<style>
  .menu-account{
    margin: 0 15px;
  }
  .avatar {
    display: inline-block;
    width: 40px;
    height: 40px;
    background: #1e67be;
    color: #fff;
    text-align: center;
    line-height: 40px;
    border-radius: 50%;
    font-size: 18px;
    text-transform: uppercase;
    cursor: pointer;
  }

  /* The container <div> - needed to position the dropdown content */

  /* Dropdown Content (Hidden by Default) */
  .dropdown-content {
    display: block;
    margin-top: 50px;
    position: absolute;
    background-color: #f9f9f9;
    min-width: 160px;
    box-shadow: 0px 8px 16px 0px rgba(0,0,0,0.2);
    z-index: 1;

    padding-left: 20px;
    right: -15px;
  }
  .menu-account .dropdown-content>* {
    vertical-align: middle;
  }
  /* Links inside the dropdown */
  .dropdown-content a {
    color: black;
    padding: 12px 0px;
    text-decoration: none;

  }

  /* Change color of dropdown links on hover */
  .dropdown-content a:hover {background-color: #f1f1f1}


  .dropdown-content i {
    padding: 0.4em;
    -webkit-transition: .1s ease-in-out all;
    transition: .1s ease-in-out all;
    border-radius: 50%;
  }

</style>
