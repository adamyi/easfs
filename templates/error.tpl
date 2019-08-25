<!DOCTYPE html>
<html lang=en>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,height=device-height,initial-scale=1.0">
  <title>Geegle ÜberProxy</title>
  <link href="https://fonts.googleapis.com/css?family=Roboto" rel="stylesheet">
  <style>
    :root {
      --color-grey-300: #dadce0;
      --color-grey-500: #9aa0a6;
      --color-grey-600: #80868b;
      --color-grey-700: #5f6368;
      --color-grey-800: #3c4043;
      --color-grey-900: #202124;
      --color-blue-50: #e8f0fe;
      --color-blue-100: #d2e3fc;
      --color-blue-200: #aecbfa;
      --color-blue-300: #84b4f8;
      --color-blue-500: #4285f4;
      --color-blue-600: #1a73e8;
      --color-blue-700: #1967d2;
    }
    @media (max-width: 700px) {
      :root {
        --card-width: 100%;
      }
      #card-container {
        position: absolute;
        left: 0;
        top: 0;
        width: 100%;
        height: 100%;
      }
      #card {
        display: flow-root;
      }
      .card-content {
        padding-top: 32px;
      }
    }
    @media (min-width: 701px) and (max-height: 700px) {
      #card-container {
        top: 20px !important;
        transform: translate(-50%, 0%) !important;
      }
    }
    @media (min-width: 701px) {
      :root {
        --card-width: 640px;
      }
      #card {
        border-radius: 8px;
        box-shadow: 0 1px 2px 0px rgba(60, 64, 67, 0.3),
                    0 1px 3px 1px rgba(60, 64, 67, 0.15);
      }
      #card-container {
        position: absolute;
        left: 50%;
        top: 50%;
        transform: translate(-50%, -50%);
      }
      .card-content {
        padding: 48px 40px 48px 40px;
      }
    }
    ::-webkit-scrollbar {
      width: 16px;
    }
    ::-webkit-scrollbar-thumb {
      border-radius: 16px;
      border: solid 4px rgba(0, 0, 0, 0);
      background: var(--color-grey-300);
      background-clip: padding-box;
    }
    html {
      font: 14px/22px 'Roboto', arial;
    }
    #card {
      width: var(--card-width);
      height: var(--card-height);
      color: var(--color-grey-700);
    }
    .card-content {
      height: calc(100% - 100px);
    }
    .brand {
      font-weight: normal;
      color: var(--color-grey-700);
      font-size: 24px;
      line-height: 22px;
    }
    .title {
      color: var(--color-grey-800);
      margin: 24px 0px 8px 0px;
      font-size: 24px;
      line-height: 32px;
      font-weight: normal;
    }
    .description {
      color: var(--color-grey-800);
      margin: 0px 0px 24px 0px;
      font-size: 18px;
      line-height: 24px;
      font-weight: normal;
    }
    a {
      color: var(--color-blue-600);
      height: 48px;
      border-radius: 3px;
      text-decoration: none;
    }
    .center {
      text-align: center;
    }
  </style>
</head>
<body>
  <div id="card-container">
    <div id="card">
      <div class="card-content">
        <div class="center">
          <span class="brand">EASFS</span>
        </div>
        <div class="center">
          <h1 class="title">{{ .Title }}</h1>
          <h3 class="description">{{ .Description }}</h3>
        </div>
      </div>
    </div>
  </div>
</body>
</html>
