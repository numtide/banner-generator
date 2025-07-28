# Font Strategy Comparison

## Current Approach: Embedded Base64 TTF

**Pros:**
- Self-contained SVG (works everywhere)
- No external dependencies
- Works offline

**Cons:**
- Large file size (~113KB per SVG)
- Base64 encoding adds 33% overhead
- Repeated font data in every banner

## Recommended Approach: External Web Fonts

**Pros:**
- Tiny SVG files (~3-5KB)
- Font cached by browser
- Better performance for multiple banners
- Uses modern WOFF2 format

**Cons:**
- Requires font hosting
- May not display custom font if server is down
- CORS considerations

## Implementation Options:

### 1. Hybrid Approach
```bash
# For embedded fonts (current, self-contained)
./bin/banner-api -template templates/banner.svg.mustache

# For web fonts (smaller, faster)
./bin/banner-api -template templates/banner-webfont.svg.mustache
```

### 2. Font Serving Setup
The API now serves fonts at `/fonts/` endpoint. Copy your web fonts:

```bash
# Create optimized font directory
mkdir -p fonts/web
cp "fonts/GT Pressura Regular/Web Fonts/2ea3b87b231d6d46708c756a0e04610e.woff2" fonts/web/gt-pressura-regular.woff2
cp "fonts/GT Pressura Regular/Web Fonts/2ea3b87b231d6d46708c756a0e04610e.woff" fonts/web/gt-pressura-regular.woff
```

### 3. CDN Approach
For production, consider serving fonts from a CDN:

```mustache
@font-face {
  font-family: 'GT Pressura';
  src: url('https://cdn.example.com/fonts/gt-pressura-regular.woff2') format('woff2'),
       url('https://cdn.example.com/fonts/gt-pressura-regular.woff') format('woff');
}
```

## File Size Comparison:

| Approach | SVG Size | Font Loading | Total First Load |
|----------|----------|--------------|------------------|
| Embedded TTF | ~113KB | Included | 113KB |
| Web Font | ~4KB | ~15KB (WOFF2) | 19KB |
| CDN Font | ~4KB | ~15KB (cached) | 4-19KB |

## Recommendations:

1. **For README embeds**: Use embedded fonts (current approach) for maximum compatibility
2. **For web dashboards**: Use external web fonts for better performance
3. **For high-traffic**: Use CDN-hosted fonts with long cache headers

## Browser Support:

- **WOFF2**: Chrome 36+, Firefox 39+, Safari 12+
- **WOFF**: Chrome 6+, Firefox 3.6+, Safari 5.1+
- **TTF**: All browsers (fallback)