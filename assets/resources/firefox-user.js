// =================
// GENERAL / STARTUP
// =================

// Client Side Decorations (CSD)
user_pref("browser.tabs.drawInTitlebar", true);

// New Window
// user_pref("browser.startup.homepage", "about:home");
// New tab
// user_pref("browser.newtabpage.enabled", true);

// Startup
// 0 Start with a blank page (about:blank).
// 1 Start with the web page(s) defined as the home page(s).
// 2 Load the last visited page.
// 3 Resume the previous browser session
user_pref("browser.startup.page", 1);
// Warm on quit
user_pref("browser.sessionstore.warnOnQuit", true);
// Check that Firefox is the default browser
user_pref("browser.shell.checkDefaultBrowser", true);


// ==========================
// FIREFOX NEWTABPAGE CONTENT
// ==========================

// Searches
user_pref("browser.newtabpage.activity-stream.showSearch", true);
// Most visited websites
user_pref("browser.newtabpage.activity-stream.feeds.topsites", true);
// Highlights (key elements)
user_pref("browser.newtabpage.activity-stream.feeds.section.highlights", false);
// News feed
user_pref("browser.newtabpage.activity-stream.feeds.snippets", false); 


// ==============
// SEARCH ENGINES
// ==============

// Moteurs de recherche accessibles en un clic désactivés
user_pref("browser.search.hiddenOneOffs", "Google,Bing,Amazon.fr,eBay");


// ========================
// GENERAL PRIVACY SETTINGS
// ========================

user_pref("browser.contentblocking.category", "strict");

// For manual configuration
// user_pref("network.cookie.cookieBehavior", 0);
// user_pref("privacy.trackingprotection.enabled", true);
// user_pref("privacy.trackingprotection.socialtracking.enabled", true);
// user_pref("privacy.trackingprotection.pbmode.enabled", false);
// user_pref("privacy.trackingprotection.enabled", true);
// user_pref("privacy.trackingprotection.pbmode.enabled", false);
// user_pref("privacy.trackingprotection.socialtracking.enabled", true);
// user_pref("privacy.trackingprotection.cryptomining.enabled", false);
// user_pref("privacy.trackingprotection.enabled", true);
// user_pref("privacy.trackingprotection.fingerprinting.enabled", true);
// user_pref("privacy.trackingprotection.socialtracking.enabled", true);


// =======================================
// NETWORK PROTOCOLS FOR IMPROVING PRIVACY
// =======================================

// In order to check the config
// https://www.cloudflare.com/ssl/encrypted-sni/

// Trusted Recursive Resolver (TRR) aka DoH
// see https://wiki.mozilla.org/Trusted_Recursive_Resolver
user_pref("network.trr.mode", 2);
user_pref("network.trr.uri", "https://mozilla.cloudflare-dns.com/dns-query");
user_pref("network.trr.bootstrapAddress", "1.1.1.1");

// Encrypted Server Name Indication
user_pref("network.security.esni.enabled", true);

// First-Party Isolation (FPI) - BREAKS TOO MANY WEBSITES end 2019
// FPI works by separating cookies on a per-domain basis
// global switch
// user_pref("privacy.firstparty.isolate", true);
// to block postMessage across different first party domains
// user_pref("privacy.firstparty.isolate.block_post_message", false);
// Users can set this parameter to false if they're having problems logging into websites (lower some of the "isolation" rules).
// user_pref("privacy.firstparty.isolate.restrict_opener_access", true);

