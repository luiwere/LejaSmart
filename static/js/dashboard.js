/**
 * LejaSmart Landing Page Interactive Scripts
 */

document.addEventListener('DOMContentLoaded', function () {
  
  // 1. Sticky Header Shadow on Scroll
  const header = document.getElementById('site-header');
  window.addEventListener('scroll', function () {
    if (window.scrollY > 20) {
      header.style.boxShadow = '0 4px 20px rgba(0, 0, 0, 0.08)';
      header.style.background = 'rgba(255, 255, 255, 0.98)';
    } else {
      header.style.boxShadow = 'none';
      header.style.background = 'rgba(255, 255, 255, 0.92)';
    }
  });

  // 2. Mobile Menu Toggle
  const mobileBtn = document.getElementById('mobile-menu-btn');
  const mobileNav = document.getElementById('mobile-nav');
  if (mobileBtn && mobileNav) {
    mobileBtn.addEventListener('click', function () {
      const isOpen = mobileNav.classList.toggle('open');
      mobileBtn.setAttribute('aria-expanded', isOpen);
    });

    // Close mobile nav when link clicked
    document.querySelectorAll('.mobile-nav-link, .mobile-nav .btn').forEach(link => {
      link.addEventListener('click', () => {
        mobileNav.classList.remove('open');
        mobileBtn.setAttribute('aria-expanded', 'false');
      });
    });
  }

  // 3. Animated Number Counters (Intersection Observer)
  const nums = document.querySelectorAll('.num');
  let animated = false;

  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting && !animated) {
        animated = true;
        nums.forEach(span => {
          const target = Number(span.getAttribute('data-target')) || 0;
          let current = 0;
          const duration = 1500;
          const startTime = performance.now();

          const updateCount = (currentTime) => {
            const elapsedTime = currentTime - startTime;
            const progress = Math.min(elapsedTime / duration, 1);
            // Ease out quad
            const easeProgress = 1 - (1 - progress) * (1 - progress);
            current = Math.floor(easeProgress * target);
            
            span.textContent = current.toLocaleString();

            if (progress < 1) {
              requestAnimationFrame(updateCount);
            } else {
              span.textContent = target.toLocaleString();
            }
          };

          requestAnimationFrame(updateCount);
        });
      }
    });
  }, { threshold: 0.2 });

  const statsSection = document.querySelector('.stats-strip');
  if (statsSection) {
    observer.observe(statsSection);
  }

  // 4. Hero Mockup Tab Switcher
  const mockupTabs = document.querySelectorAll('.mockup-tab');
  const mockupPanels = document.querySelectorAll('.mockup-tab-panel');

  mockupTabs.forEach(tab => {
    tab.addEventListener('click', function () {
      const targetId = this.getAttribute('data-mockup-target');
      
      mockupTabs.forEach(t => t.classList.remove('active'));
      mockupPanels.forEach(p => p.classList.remove('active'));

      this.classList.add('active');
      const activePanel = document.getElementById(targetId);
      if (activePanel) {
        activePanel.classList.add('active');
      }
    });
  });

  // 5. Interactive Voice Simulator Demo in Hero Mockup
  const voiceBtn = document.getElementById('voice-demo-btn');
  const voiceTranscript = document.getElementById('voice-transcript-text');
  
  const sampleVoicePhrases = [
    {
      phrase: `"Paid KES 1,200 for Boda Boda transport to town market and packing supplies"`,
      amount: "KES 1,200",
      category: "Transport",
      note: "Boda Boda to town market"
    },
    {
      phrase: `"Bought 10 bags of Maize flour from Unga Ltd for KES 18,500 cash"`,
      amount: "KES 18,500",
      category: "Supplies",
      note: "Unga Ltd Restock"
    },
    {
      phrase: `"Paid electricity and Kenya Power token bill KES 3,400"`,
      amount: "KES 3,400",
      category: "Utilities",
      note: "KPLC Electricity token"
    },
    {
      phrase: `"Casual worker daily wages for offloading sacks KES 800"`,
      amount: "KES 800",
      category: "Other",
      note: "Casual offloading wages"
    }
  ];

  let sampleIdx = 0;
  if (voiceBtn && voiceTranscript) {
    voiceBtn.addEventListener('click', function () {
      voiceBtn.style.transform = 'scale(1.25)';
      voiceTranscript.textContent = "🎙 Listening in real-time...";
      voiceTranscript.style.color = "#2563eb";

      setTimeout(() => {
        sampleIdx = (sampleIdx + 1) % sampleVoicePhrases.length;
        const currentSample = sampleVoicePhrases[sampleIdx];
        voiceTranscript.textContent = currentSample.phrase;
        voiceTranscript.style.color = "#0f172a";
        voiceBtn.style.transform = 'scale(1)';

        const tagsContainer = document.querySelector('.voice-parsed-tags');
        if (tagsContainer) {
          tagsContainer.innerHTML = `
            <span class="tag">Amount: <strong>${currentSample.amount}</strong></span>
            <span class="tag">Category: <strong>${currentSample.category}</strong></span>
            <span class="tag">Status: <strong>Auto-Saved to Expenses ✓</strong></span>
          `;
        }
      }, 700);
    });
  }

  // 6. Role Solutions Tab Switcher
  const roleTabs = document.querySelectorAll('.role-tab');
  const roleContents = document.querySelectorAll('.role-tab-content');

  roleTabs.forEach(tab => {
    tab.addEventListener('click', function () {
      const role = this.getAttribute('data-role');

      roleTabs.forEach(t => t.classList.remove('active'));
      roleContents.forEach(c => c.classList.remove('active'));

      this.classList.add('active');
      const targetContent = document.getElementById(`role-${role}`);
      if (targetContent) {
        targetContent.classList.add('active');
      }
    });
  });

  // 7. FAQ Accordion
  const faqQuestions = document.querySelectorAll('.faq-question');
  faqQuestions.forEach(btn => {
    btn.addEventListener('click', function () {
      const item = this.parentElement;
      const isOpen = item.classList.contains('open');

      // Close other open items
      document.querySelectorAll('.faq-item').forEach(i => {
        if (i !== item) {
          i.classList.remove('open');
          i.querySelector('.faq-question').setAttribute('aria-expanded', 'false');
        }
      });

      // Toggle clicked item
      item.classList.toggle('open', !isOpen);
      this.setAttribute('aria-expanded', !isOpen);
    });
  });

});

