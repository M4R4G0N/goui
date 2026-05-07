package components

import (
	"fmt"
	"strings"
)

type TabItem struct {
	Text     string
	Href     string
	Active   bool
	TargetID string
}

type TabsComponent struct {
	ID    string
	Items []TabItem
}

func (t *TabsComponent) Render() string {
	if t.ID == "" {
		t.ID = AutoID()
	}

	targets := []string{}
	for _, item := range t.Items {
		if item.TargetID != "" {
			targets = append(targets, item.TargetID)
		}
	}
	targetsJS := fmt.Sprintf("['%s']", strings.Join(targets, "','"))

	var tabsHTML strings.Builder
	for _, item := range t.Items {
		activeClass := ""
		if item.Active {
			activeClass = "active"
		}
		href := item.Href
		if href == "" {
			href = "javascript:void(0)"
		}
		tabsHTML.WriteString(fmt.Sprintf(
			`<a href="%s" class="goui-tab %s" data-target="%s" onclick="switchTab(this,'%s','%s',%s)">%s</a>`,
			href, activeClass, item.TargetID, t.ID, item.TargetID, targetsJS, item.Text,
		))
	}

	script := `<script>
if (typeof window.switchTab !== 'function') {
  window.switchTab = function(el, groupId, targetId, allTargets) {
    var container = el.closest('.goui-tabs');
    container.querySelectorAll('.goui-tab').forEach(function(t){ t.classList.remove('active'); });
    el.classList.add('active');
    if (targetId && allTargets) {
      allTargets.forEach(function(id){
        var t = document.getElementById(id);
        if (t) t.style.display = 'none';
      });
      var target = document.getElementById(targetId);
      if (target) target.style.display = 'block';
    }
    if (history.replaceState) history.replaceState(null, null, '#' + targetId);
  };
}
if (typeof window.gouiActivateHash !== 'function') {
  window.gouiActivateHash = function() {
    var hash = window.location.hash.replace('#', '');
    if (!hash) return;
    var tab = document.querySelector('.goui-tab[data-target="' + hash + '"]');
    if (!tab) return;
    var container = tab.closest('.goui-tabs');
    if (!container) return;
    var targets = [];
    container.querySelectorAll('.goui-tab[data-target]').forEach(function(t) {
      targets.push(t.getAttribute('data-target'));
    });
    window.switchTab(tab, container.id, hash, targets);
  };
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', window.gouiActivateHash);
  } else {
    window.gouiActivateHash();
  }
  window.addEventListener('hashchange', window.gouiActivateHash);
}
</script>`

	return fmt.Sprintf(`<div class="goui-tabs" id="%s">%s</div>%s`, t.ID, tabsHTML.String(), script)
}

func Tabs(items ...TabItem) *TabsComponent {
	return &TabsComponent{Items: items}
}
