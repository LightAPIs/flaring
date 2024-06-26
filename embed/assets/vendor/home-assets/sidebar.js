(() => {
  const hideSidebar = 'hide-sidebar';
  const showSidebar = 'show-sidebar';
  const categoryGroupContainer = document.querySelector('.category-group-container');
  if (categoryGroupContainer) {
    let state = true;
    const doc = document.documentElement;
    if (doc.clientWidth <= 768) {
      state = false;
    }

    function setSidebarState(sta) {
      if (sta !== undefined) {
        state = sta;
      } else {
        state = !state;
      }

      if (state) {
        categoryGroupContainer.classList.remove(hideSidebar);
        categoryGroupContainer.classList.add(showSidebar);
      } else {
        categoryGroupContainer.classList.remove(showSidebar);
        categoryGroupContainer.classList.add(hideSidebar);
      }
    }

    /** 取消搜索输入框的焦点 */
    function blurSerachInput() {
      /** 搜索框 */
      const searchInput = document.getElementById('search');
      setTimeout(() => {
        searchInput?.blur();
      }, 0);
    }

    document.querySelectorAll('.catetory-group-content .category-list a').forEach(ele => {
      const id = ele.getAttribute('href').replace('#', '');
      ele.addEventListener('click', ev => {
        ev.stopPropagation();
        ev.preventDefault();
        const categoryTitle = document.querySelector(`[data-set-category="${id}"]`);
        if (categoryTitle) {
          const offsetTop = categoryTitle.offsetTop - 8;
          doc.scrollTo({
            top: offsetTop,
            behavior: 'smooth',
          });
          const clientWidth = doc.clientWidth;
          if (clientWidth <= 768) {
            // 隐藏侧边栏
            setSidebarState(false);
          }
        } else {
          console.warn(`Unable find category title about id "${id}".`);
        }
      });
    });

    const sidebarBtn = document.getElementById('btn-category-list');
    if (sidebarBtn) {
      sidebarBtn.addEventListener('click', ev => {
        ev.stopPropagation();
        ev.preventDefault();
        blurSerachInput();
        setSidebarState();
      });
      sidebarBtn.addEventListener(
        'touchstart',
        ev => {
          ev.stopPropagation();
          ev.preventDefault();
          blurSerachInput();
          setSidebarState();
        },
        {
          passive: false,
        }
      );
    }
  }
})();
