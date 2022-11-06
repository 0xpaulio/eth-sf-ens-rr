import './App.css';

import {useState} from "react";

import {ethers} from "ethers";
import {CONTRACT_ABI, L1_REGISTRY_CONTRACT_ADDRESS} from "./contract";

function App() {
    const [connectedAccount, setConnectedAccount] = useState("");

    const nullModal = {open: false, functionSelectorName: "", response: ""}
    const [modal, setModal] = useState(nullModal);
    const [loading, setLoading] = useState("");

    const ensName = "paulio.eth";
    const txOptions = {ccipReadEnabled: true};

    const connectWallet = async () => {
        try {
            const {ethereum} = window;

            if (!ethereum) {
                alert("No wallet connect found")
                return
            }

            const accounts = await ethereum.request({method: "eth_requestAccounts"});
            setConnectedAccount(accounts[0]);
        } catch (err) {
            handleError(err)
        }
    }

    const createRegistryContract = (ethereum) => {
        return new ethers.Contract(
            L1_REGISTRY_CONTRACT_ADDRESS,
            CONTRACT_ABI,
            (new ethers.providers.Web3Provider(ethereum))
        );
    }

    const getOwner = async (ensName) => {
        try {
            const {ethereum} = window;

            if (ethereum) {
                setLoading("Owner")

                // Initialize the registry contract connection
                const registry = createRegistryContract(ethereum);

                // Generate the nameHash for the corresponding ENS name
                const node = ethers.utils.namehash(ensName);

                // Query the L1 Registry contract & let ethers handle the CCIP reversion
                // const owner = await registry.owner(node, txOptions)
                const owner = await registry.ownerWithProof("0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000002a00000000000000000000000000000000000000000000000000000000000001280000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000007290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563633dc4d7da7256660a892f8f1604a44b5432649cc8ec5cb3ced4c4e6ac94dd1d890740a8eb06ce9be422cb8da5cdafc2b58c0a5e24036c578de2a433c828ff7d3b8ec09e026fdc305365dfc94e189a81b38c7597b3d941c279f042e8206e0bd8d6f1980eb2314f67f63326c8e41b95464bf34a84ab89ae1e7d91ebf97a59b82874f0d6cace4a2b05d60805e71590d9cce5bf60ab136621e647ac2886b6c21be069a286d0f158a81c833f9c9debb57bf7a769cc002bc216d751fdd72ca0338f210000000000000000000000000000000000000000000000000000000000000fb5f90fb2b9042a30786639303231316130646239633964616132303533366232323665366162383833663065663534343562323862663032396261613034393035333637396336646534303662356666356130616561303033636665623562356437383436313230316334316435656139316434646531396631393762343531646365356264313338376537366135323837386130333732393831666539376664373034376430643165636135626338383163346563653862643538313230656431666436623033386333303739636366363938396130396537373230633833656461323163383437313164303163333131636664363535343539303732396638333662666330386435393764663137623264326465376130383961323030623065343764343332636565353135353139383538343763323863366266323864646536323334616532303438363535343238323136343230316130636338656465663665636666343030303231346231616132313537393564643138366662656362353931303837636366383533663933353939643939336162306130393337353333363437313935643737643430363936306539323264646137646537623865326561336230313966333937366432306165346534653863646434386130636165393962323065663036316163613361313436646530636463386131613866363062383738636236393364303132386635613461336164383731643761636130373636336635356263323262333534306564313961373165333237333235353361636432633361656432633432353135366662333762333764633963323336646130346538303639383866363036636336383930373466636533326239353236353634666264316636393164376363303738303435336632366235366434386233356130363733653839633061666431636330396430306361303337343462343863623530663264303037363737363463313838306262643163366431636138373832316130656564303537623435633330626239643033363166336633383962643639613636393438333736333238386135396664613233316163623039393437616433366130663735656336323365346365343631633162366562313034633733343638333034663833343666666433656532306239333239393165623664343539333135626130626662383130326165636161336362396630396236306137383665313665663037306365353165313265343664313865323838366366343936366333613933336130653637323965366235306536303233633263333036656562623838343232326661343664393231643765346663346636636462356266623664353030646266336130613432333235656664616161656630663833613563386264643032353739656364646234326236663032643832316561646630336661353838396164383362303830b9042a30786639303231316130313165303735623032616435343037616634623864336536303636386463386263393531363666646637326136323661666266346262626434323535663931386130343131306134353836653237653165336462333537626531383666343362656236653238643466663834616430373031336239656132646365303338363864646130353039653936376632383038623935396138376131333266393137393139376631393064373930366431303162373461343235373434613132333463313366346130343861643334326635393839323330633838636464323730386638383261376536336637323365393361336330663165633431306434323836383661653530336130306462643662363836393438313430353264326637323133326661393933346462373239656535643330303233626262653039313663303136393939613338346130343633336632663332663037646532376364343763666238303566383839616664343965653837396634663135386166616532633064386234653534363461616130303963356137383232633264353830626464666633313838316634336133323866626630346534383630363362633466373236313061613164623131356639616130303836393464613738636137356533353838336565613333653337626536346261666661323833303264306337613835323836373562313637373266313461656130643832633134373464643966616532633539353532333030353066303335626537633663316634346233643231653266646537303362666335393839303433356130623135643038393266313463343331336161633063373366346163343365363638346163346661383936373466653831386534323861393632373661643130626130386465316636373739623339343139323434323333623934303939393966373832623564393538633061626537643739623037656436346237643333346233356130356261626131363766306366316564653839356436383635373832313934383730663866323838666462313665323265653562626264666165363234646532656130336462383836626266386139376330343930643635313463323838623338666437396639613761373233303161653939353366386564393731616666343236376130643264366637346231666536333632656632613431653336643335663062353035616362613863636631626161346463656333663766336261373336376134656130623163363364656463396363313161646238383733336530396337636331373638356265343637333832393439353364663730626538306364343363386463306130393735383733636234386535363262356335356232373866653864666362663165316333353631616237623637303930636330653531623834346637623166343830b9042a30786639303231316130386239343232323561363936616439373861616636396562393734376530356336626562353336353061376266373561393331643538653464333234386261346130363639383062353166626233653764326331323630313263636663363832343130363466363131343665326466376432643236333964383736303664316637346130336362653565353338313464613934333135346435333265343662376135333864326339383034623265363636396634643133373336653262653061623362396130343263666265383335303064616262343331613061396234623730326162666334383462393430636136646136313234363732303838633536346631383663666130396663333939333232376134393735663161383839393939346231613462346630313436623037306564316639633164303036366664396534313865643264386130646431396233343863656363313231353766656435336235663265326463653366316664623963336438366134393836336132353535623239616438333962646130383564376535633263353835643439303032373365653038343137313438633334323761316630633733303165663035633232396562343839306136663333646130626438363238336165316439643661343636316639343563633361656439643131353463633463633234393635366431306631386261343633653434313663336130613631326561656463633636353638633965313865663431393030336466636165366265636637663564316264393334623161623261356262383138366332646130633733653430313938393134393238366266653066386232316430646239363539383334376139383939326438623938643534336531653262373665636334646130643439363839366138313630623330333431363534646439393863313661353266376162653331666430623339666233343762343562363433653530353938636130663935653861336339303539373636636366666333343663336166343330356433633138393435393635366362343737646266373165353163613638333135366130303230323238336139656261643464396131663037396263353938666163623164663037393632343165333233306564313538663238666439633338353836656130313931393663653536326333343265666661666535613639613261633964656235633261313537386133666338356133646463626538313938616132303564326130366330313064633462303336313233616431616235663238393038363934316633653130376336353536366233356533636262316234353433663139376137626130636664353263373330343930613633646632363834326562636136366136326232356638356466626364646638386439666435643366356335646337663465313830b901a830786638643138303830383061303463323434643837373964613761613338376437663333363361616530363065316666646133633566396562633239613666623437393338616338326434326661303536343432616236356131666461633332363130386264626162643865653337666530643565636664346561356661613063636566396233366432336232323938303830613030623039653631663035323761396139613862323035356238366438626631393262383733333731613262316461636161356635643933366432316564306165383038303830383061303439626635613964306461616366353962376239363862666631316139656232656535313334333965666665373463393265396661656135313864376664613961303436366366623233353162356138323339313933323365633834643237643034383534353439306130663738646133666335386430366166316530663537643861303737353036306632383465393935323639313139363738616664316165393333343230343538623963616137373664346435386430383134376635383438633338303830b8a8307866383531383038306130303463333831623261316636653534306165333338656534333130643432643865383533303266643363653131353866623365613230633964616261313336623830383038303830383038303830383061306633666135646537306638323339646463613538313938346566356464336631623062383731663964373335633136636633616332343131626165323665343238303830383038303830b8d4307866383637396533623139306638666436396337373630393863643762306465623262396533353434616437623932333862376238633531326666373633393535623462383436663834343031383061303232326161636361653230616635333139366261353461666561663037306561323837626465396630396136633064353130646463373265616534353365653661303838336665633932343661333537613831623338316137373335616238366533396436386364626637383362626365396162336431303834323862346439343400000000000000000000000000000000000000000000000000000000000000000000000000000000000281f9027ef9027bb8423078633138323161373661396335636131646463646236333838373836643865656566343065386633626430363963333464613634313766363939323161363730309484c970bfcd59a0e98ec6f13cbdf24aa1a741f033f9021fb901a830786638643138303830613039333265616432393438306233306537383133333433363565366632666666656338666136653933366136666537616435613965343234383964313362333437383038306130373737646263346538363530363961663136336239666461353635333664346634656631303439666165613931383430396663666337346633396362653136616130303464306632623065353039666466393366343832613534663335343066613030356464636631383362353636343031326239636133626663323537356263323830383061303062636334366231313836393365363563636436303261356136346339393738623463656131386336653432376434653564656166626639356238666362336238306130653338303436343462646336303331393736386330633832616530613234623030336539633532306361316666393466383535633965613362613730353964373830383038306130333363356235663138633738636434626235333363353836376165656136633233333833333538663439373962356232653763663933346537346638633031333830b87230786637613033316364646637373739323333373139623838333366653562613136616330306364323236643966306139393061343961303065336566386261383136666663393539343834633937306266636435396130653938656336663133636264663234616131613734316630333300000000000000000000000000000000000000000000000000000000000000", node)
                setModal(
                    {
                        open: true,
                        functionSelectorName: "Owner()",
                        response: owner
                    }
                )
                setLoading("")
                return owner
            }
        } catch (err) {
            handleError(err)
        }
    }

    const getResolver = async () => {
        try {
            const {ethereum} = window;

            if (ethereum) {
                setLoading("Resolver")

                // Initialize the registry contract connection
                const registry = createRegistryContract(ethereum);

                // Generate the nameHash for the corresponding ENS name
                const node = ethers.utils.namehash(ensName);

                // Query the L1 Registry contract & let ethers handle the CCIP reversion
                const resolver = await registry.resolver(node, txOptions)
                setModal(
                    {
                        open: true,
                        functionSelectorName: "Resolver()",
                        response: resolver
                    }
                )
                setLoading("")
                return resolver
            }
        } catch (err) {
            handleError(err)
        }
    }

    const getTTL = async () => {
        try {
            const {ethereum} = window;

            if (ethereum) {
                setLoading("TTL")

                // Initialize the registry contract connection
                const registry = createRegistryContract(ethereum);

                // Generate the nameHash for the corresponding ENS name
                const node = ethers.utils.namehash(ensName);

                // Query the L1 Registry contract & let ethers handle the CCIP reversion
                const ttl = await registry.ttl(node, txOptions)
                setModal(
                    {
                        open: true,
                        functionSelectorName: "TTL()",
                        response: ttl
                    }
                )
                setLoading("")
                return ttl
            }
        } catch (err) {
            handleError(err)
        }
    }
    const handleError = (err) => {
        console.log("Handling error")
        console.log(err);
        console.log("opening modal")
        setModal(
            {
                open: true,
                functionSelectorName: "<Error/>",
                response: `${err}`
            }
        )
        setLoading("")
    }

    const loadingButton = (name, handler) => {
        if (loading === name) {
            return (
                <div className="button">
                    <div className="spin"/>
                </div>
            )
        }
        return (
            <div
                className="button"
                onClick={() => handler(ensName)}
            >
                Get {name}
            </div>
        )
    }

    return (
        <div className="App">
            {
                modal.open &&
                <div className="modal-box">
                    <div className="modal-container">
                        <h1>{modal.functionSelectorName}</h1>
                        <div className="left-align">
                            <h2 className="row">ENS Name: <p className="data">{ensName}</p></h2>
                            <h2 className="row">Response: <p className="data">{modal.response}</p></h2>
                        </div>
                        <div
                            className="button center"
                            onClick={() => setModal(nullModal)}
                        >
                            Close Modal
                        </div>
                    </div>
                </div>
            }

            {/* Header */}
            <div className="Header">
                <h1>Trustless L2 ENS Registry Rollup Solution</h1>

                {
                    connectedAccount === "" ?
                        <div
                            className="button connect-wallet"
                            onClick={connectWallet}
                        >
                            Connect Wallet
                        </div>
                        :
                        <div
                            className="button connect-wallet disabled"
                        >
                            {`${connectedAccount.substr(0, 6)}...${connectedAccount.substr(connectedAccount.length - 4)}`}
                        </div>
                }
            </div>

            {/* Body */}
            <div className="Body">
                <h2>Cross Chain Accessor Functions:</h2>
                {loadingButton("Owner", getOwner)}
                {loadingButton("Resolver", getResolver)}
                {loadingButton("TTL", getTTL)}
            </div>

            {/* Footer */}
            <div className="Footer">
                <h1>Made with ❤</h1>
            </div>
        </div>
    );
}

export default App;
